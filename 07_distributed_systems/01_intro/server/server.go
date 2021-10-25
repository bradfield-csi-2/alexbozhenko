package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	SERVER           = "127.0.0.1"
	PORT             = "8000"
	URL              = SERVER + ":" + PORT
	STORAGE_FILEPATH = "storage"
)

type inMemoryStorage map[string]string

var keyValueMap = make(inMemoryStorage)

func getHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["key"]
	if !ok {
		fmt.Fprintf(w, "Key parameter is not found in the request")
		return
	}
	key := keys[0]
	value, ok := keyValueMap[key]
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

	f, err := os.OpenFile(STORAGE_FILEPATH, os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	//	b := new(bytes.Buffer)
	encoder := gob.NewEncoder(f)
	err = encoder.Encode(*m)
	return err
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	reqData := struct {
		Key   string
		Value string
	}{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(&w, "Bad request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		errorResponse(&w, "Bad request", http.StatusBadRequest)
		return
	}
	_, keyExists := keyValueMap[reqData.Key]
	keyValueMap[reqData.Key] = reqData.Value
	err = persistUpdate(&keyValueMap)
	if err != nil {
		errorResponse(&w, "Server error", http.StatusInternalServerError)
		return
	}
	if keyExists {
		w.Write([]byte("UPDATED"))
	} else {
		w.Write([]byte("INSERTED"))
	}
}

func init() {
	// data, err 	:= os.ReadFile(STORAGE_FILEPATH)
	f, err := os.OpenFile(STORAGE_FILEPATH, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic("Error reading storage")
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&keyValueMap)
	if err != nil {
		panic("Error decoding storage")
	}
}

func main() {
	fmt.Println("Welcome to the distributed K-V store server")
	//TODO init

	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/put", putHandler)

	http.ListenAndServe(URL, nil)
}
