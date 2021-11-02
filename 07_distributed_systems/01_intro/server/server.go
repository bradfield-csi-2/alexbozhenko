package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"kv_store/protocol"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
)

const (
	SERVER           = "127.0.0.1"
	PORT             = "8000"
	URL              = SERVER + ":" + PORT
	STORAGE_FILEPATH = "storage"
)

// abozhenko for oz: Ok, now we have this map that fits into memory
// Should we assume that whole dataset does not fit in memory, and
// even on disk on a single server?
type inMemoryStorage map[string]string

var keyValueMap = make(inMemoryStorage)

var storageFD *os.File
var mutex sync.RWMutex

func getHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getHandler. Goroutines:%v", runtime.NumGoroutine())
	reqData := protocol.GetRequest{}
	body, err := ioutil.ReadAll(r.Body)
	// abozhenko for oz: Is there a best practice to handle errors
	// and avoid duplication in golang?
	if err != nil {
		errorResponse(&w, "Bad request", http.StatusBadRequest)
		return
	}
	err = reqData.Decode(body)
	if err != nil {
		errorResponse(&w, "Bad request", http.StatusBadRequest)
		return
	}

	key := reqData.Key
	mutex.RLock()
	value, ok := keyValueMap[string(key)]
	mutex.RUnlock()
	if ok {
		fmt.Fprintf(w, "%s", value)
	} else {
		fmt.Fprintf(w, "key %s NOT FOUND", key)
	}
}

func errorResponse(w *http.ResponseWriter, message string, code int) {
	(*w).WriteHeader(http.StatusBadRequest)
	(*w).Write([]byte(message))
}

func persistUpdate(m *inMemoryStorage) error {
	storageFD.Truncate(0)
	storageFD.Seek(0, io.SeekStart)
	encoder := gob.NewEncoder(storageFD)
	err := encoder.Encode(*m)
	return err
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("putHander. Goroutines:%v", runtime.NumGoroutine())
	reqData := protocol.SetRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(&w, "Bad request", http.StatusBadRequest)
		return
	}
	err = reqData.Decode(body)
	if err != nil {
		errorResponse(&w, "Bad request", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	_, keyExists := keyValueMap[string(reqData.Key)]
	keyValueMap[string(reqData.Key)] = string(reqData.Value)
	err = persistUpdate(&keyValueMap)
	mutex.Unlock()
	if err != nil {
		message := fmt.Sprintf("Server error: %s", err)
		errorResponse(&w, message, http.StatusInternalServerError)
		return
	}
	if keyExists {
		w.Write([]byte("UPDATED"))
	} else {
		w.Write([]byte("INSERTED"))
	}
}

func init() {
	f, err := os.OpenFile(STORAGE_FILEPATH, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error reading storage %s", STORAGE_FILEPATH))
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&keyValueMap)
	if err != io.EOF && err != nil {
		panic(fmt.Sprintf("Error decoding storage %s: %s", STORAGE_FILEPATH, err))
	}
}

func main() {
	fmt.Println("Welcome to the distributed K-V store server")
	var err error
	storageFD, err = os.OpenFile(STORAGE_FILEPATH, os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		panic(err)
	}
	defer storageFD.Close()
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/put", putHandler)

	http.ListenAndServe(URL, nil)
}
