package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type walRecord struct {
	id         uint64
	key, value string
}

func addToWAL(server *serverState, key, value string, w http.ResponseWriter) error {
	wr := walRecord{
		id:    server.currentWALId + 1,
		key:   key,
		value: value,
	}
	err := server.walGobEncoder.Encode(&wr)
	if errorResponse(&w, err, http.StatusInternalServerError) != nil {
		server.mutex.Unlock()
		return err
	}
	// incrementing here to avoid inconsistency
	server.currentWALId++
	return nil
}

func readWAL(server *serverState, minId uint64) (inMemoryStorage, error) {
	walFD, err := os.OpenFile(WAL_FILEPATH, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return inMemoryStorage{}, err
	}
	decoder := gob.NewDecoder(walFD)
	var wal []walRecord
	err = decoder.Decode(&wal)
	if err != io.EOF && err != nil {
		return inMemoryStorage{}, err
	}
	diffStorage := inMemoryStorage{
		transactionID: 0,
		kv:            make(map[string]string),
	}
	var record walRecord
	for _, record = range wal {
		if record.id < minId {
			continue
		}
		diffStorage.kv[record.key] = record.value
		diffStorage.transactionID = record.id
	}
	return diffStorage, nil
}

func (server *serverState) walGetHandler(w http.ResponseWriter, r *http.Request) {
	/*	req, err := http.NewRequest(http.MethodGet, URL+"/get", nil)
		must(err)
		q := req.URL.Query()
		q.Add("key", key)
		req.URL.RawQuery = q.Encode()
	*/

	walNumbers, ok := r.URL.Query()["since"]
	if !ok {
		errorResponse(&w,
			fmt.Errorf("since parameter is not found in the request"),
			http.StatusBadRequest)
		return
	}

	walNumber, err := strconv.ParseUint(walNumbers[0], 10, 64)
	if errorResponse(&w, err, http.StatusBadRequest) != nil {
		return
	}
	diff, err := readWAL(server, walNumber)
	if errorResponse(&w, err, http.StatusInternalServerError) != nil {
		return
	}
	encoder := gob.NewEncoder(w)
	encoder.Encode(diff)
}
