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
	Id         uint64
	Key, Value string
}

func addToWAL(server *serverState, key, value string, w http.ResponseWriter) error {
	wr := walRecord{
		Id:    server.currentWALId + 1,
		Key:   key,
		Value: value,
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
		TransactionID: 0,
		Kv:            make(map[string]string),
	}
	var record walRecord
	for _, record = range wal {
		if record.Id < minId {
			continue
		}
		diffStorage.Kv[record.Key] = record.Value
		diffStorage.TransactionID = record.Id
	}
	return diffStorage, nil
}

func (server *serverState) walGetHandler(w http.ResponseWriter, r *http.Request) {
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
