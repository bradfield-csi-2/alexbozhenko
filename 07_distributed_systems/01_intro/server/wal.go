package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	var wal walRecord
	if err != io.EOF && err != nil {
		return inMemoryStorage{}, err
	}
	diffStorage := inMemoryStorage{
		TransactionID: server.currentWALId,
		Kv:            make(map[string]string),
	}

	for err = nil; err != io.EOF; {
		err = decoder.Decode(&wal)
		if err != io.EOF && err != nil {
			continue
		}
		if wal.Id < minId {
			continue
		}
		diffStorage.Kv[wal.Key] = wal.Value
		diffStorage.TransactionID = wal.Id
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

func (server *serverState) pullPrimaryForUpdatesForever() {
	for {
		server.pullPrimaryForUpdates()
		time.Sleep(1000 * time.Millisecond)
	}

}

func (server *serverState) pullPrimaryForUpdates() {
	req, err := http.NewRequest(http.MethodGet, "http://"+PRIMARY_URL+"/async-catchup", nil)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	q := req.URL.Query()
	q.Add("since", fmt.Sprint((server.inMemoryStorage.TransactionID + 1)))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer resp.Body.Close()
	//	body, err := io.ReadAll(resp.Body)
	//fmt.Printf("%s", hex.Dump(body))
	//return
	decoder := gob.NewDecoder(resp.Body)
	catchUpRecords := inMemoryStorage{
		TransactionID: 0,
		Kv:            map[string]string{},
	}
	err = decoder.Decode(&catchUpRecords)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	server.mutex.Lock()
	for key, value := range catchUpRecords.Kv {
		server.inMemoryStorage.Kv[key] = value
	}
	err = persistUpdateWithLockedMutex(server)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return

	}
	oldRecordId := server.inMemoryStorage.TransactionID
	server.inMemoryStorage.TransactionID = catchUpRecords.TransactionID
	server.mutex.Unlock()
	log.Printf("Processed %d updates: %d to %d\n",
		len(catchUpRecords.Kv),
		oldRecordId,
		catchUpRecords.TransactionID)

}
