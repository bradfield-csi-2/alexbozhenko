package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	SERVER           = "127.0.0.1"
	PORT             = "8000"
	URL              = SERVER + ":" + PORT
	STORAGE_FILEPATH = "storage"
)

var keyValueMap = make(map[string]string)

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
	if keyExists {
		w.Write([]byte("UPDATED"))
	} else {
		w.Write([]byte("INSERTED"))
	}
}

// func init() {
// 	b := new(bytes.Buffer)
// 	decoder = gob.NewDecoder(b)
// 	err := decoder.Decode(keyValueMap)
// 	f, err := ioutil.ReadFile(STORAGE_FILEPATH)
// 	if err != nil {
// 		panic()
// 	}
// 	b.Write(f)
// 	defer close(f)

// }

func main() {
	fmt.Println("Welcome to the distributed K-V store server")
	//TODO init

	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/put", putHandler)

	http.ListenAndServe(URL, nil)
}
