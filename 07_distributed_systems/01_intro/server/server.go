package main

import (
	"bytes"
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

var mode serverMode

const (
	SERVER                     = "127.0.0.1"
	PRIMARY_PORT               = "8000"
	SYNCHRONOUS_FOLLOWER_PORT  = "8001"
	ASYNCHRONOUS_FOLLOWER_PORT = "8002"
	PRIMARY_URL                = SERVER + ":" + PRIMARY_PORT
	SYNC_FOLLOWER_URL          = SERVER + ":" + SYNCHRONOUS_FOLLOWER_PORT
	ASYNC_FOLLOWER_URL         = SERVER + ":" + ASYNCHRONOUS_FOLLOWER_PORT
	STORAGE_FILE_PREFIX        = "storage_"
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
	if errorResponse(&w, err, http.StatusBadRequest) != nil {
		return
	}
	err = reqData.Decode(body)
	if errorResponse(&w, err, http.StatusBadRequest) != nil {
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

func errorResponse(w *http.ResponseWriter, err error, code int) error {
	if err != nil {
		message := fmt.Sprintf("%s: %s", http.StatusText(code), err)
		(*w).WriteHeader(code)
		(*w).Write([]byte(message))
		return err
	}
	return nil
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
	if errorResponse(&w, err, http.StatusBadRequest) != nil {
		return
	}
	err = reqData.Decode(body)
	if errorResponse(&w, err, http.StatusBadRequest) != nil {
		return
	}
	mutex.Lock()
	_, keyExists := keyValueMap[string(reqData.Key)]
	if mode == PRIMARY {
		//send updates to the synchronous follower
		reqBytes := reqData.Encode()
		req, err := http.NewRequest(http.MethodPut, "http://"+SYNC_FOLLOWER_URL+"/put",
			bytes.NewReader(reqBytes))
		if errorResponse(&w, err, http.StatusInternalServerError) != nil {
			return
		}
		req.Header.Set("Content-Type", "application/octet-stream")
		resp, err := http.DefaultClient.Do(req)
		if errorResponse(&w, err, http.StatusInternalServerError) != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("got %d(%s) from synchornous replica. Not proceeding with the write",
				resp.StatusCode,
				http.StatusText(resp.StatusCode))
			errorResponse(&w, err, http.StatusInternalServerError)
			return
		}
	}
	keyValueMap[string(reqData.Key)] = string(reqData.Value)
	err = persistUpdate(&keyValueMap)
	mutex.Unlock()
	if errorResponse(&w, err, http.StatusInternalServerError) != nil {
		return
	}
	if keyExists {
		w.Write([]byte("UPDATED"))
	} else {
		w.Write([]byte("INSERTED"))
	}
}

func readStorage(filename string) {
	f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error reading storage %s", filename))
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&keyValueMap)
	if err != io.EOF && err != nil {
		panic(fmt.Sprintf("Error decoding storage %s: %s", filename, err))
	}
}
func usage() {
	fmt.Fprintf(os.Stderr, "usage: server <primary|sync_follower|async_follower>\n")
	os.Exit(1)
}

func main() {
	nArgs := len(os.Args)
	if nArgs != 2 {
		usage()
	}
	var url string
	switch os.Args[1] {
	case "primary":
		mode = PRIMARY
		url = PRIMARY_URL
	case "sync_follower":
		mode = SYNCHRONOUS_FOLLOWER
		url = SYNC_FOLLOWER_URL
	case "async_follower":
		mode = ASYNCHRONOUS_FOLLOWER
		url = ASYNC_FOLLOWER_URL
	default:
		usage()
	}

	filename := STORAGE_FILE_PREFIX + fmt.Sprint(mode)
	readStorage(filename)
	fmt.Printf("Welcome to the distributed K-V store server in %s mode\n", mode)
	fmt.Printf("Listening on: %s\n", url)
	fmt.Printf("Using storage file: %s\n", filename)

	var err error
	storageFD, err = os.OpenFile(filename, os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		panic(err)
	}
	defer storageFD.Close()
	http.HandleFunc("/get", getHandler)

	// Idea is that writes should be accepted only by the primary node,
	// but it is not enforced, since we do not have any auth anyway
	http.HandleFunc("/put", putHandler)

	http.ListenAndServe(url, nil)
}
