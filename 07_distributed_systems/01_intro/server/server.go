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
	"strconv"
	"sync"
)

const (
	SERVER                     = "127.0.0.1"
	PRIMARY_PORT               = 8000
	PARTITION_PORT_OFFSET      = 1000
	SYNCHRONOUS_FOLLOWER_PORT  = "8001"
	ASYNCHRONOUS_FOLLOWER_PORT = "8002"
	SYNC_FOLLOWER_URL          = SERVER + ":" + SYNCHRONOUS_FOLLOWER_PORT
	ASYNC_FOLLOWER_URL         = SERVER + ":" + ASYNCHRONOUS_FOLLOWER_PORT
	STORAGE_FILE_PREFIX        = "storage_"
	WAL_FILEPATH               = "wal"
)

type inMemoryStorage struct {
	TransactionID uint64
	Kv            map[string]string
}

type serverState struct {
	walFD           *os.File
	walGobEncoder   *gob.Encoder
	mode            serverMode
	inMemoryStorage inMemoryStorage
	mutex           sync.RWMutex
	storageFD       *os.File
}

func (server *serverState) getHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getHandler. Goroutines:%v", runtime.NumGoroutine())
	reqData := protocol.GetRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if errorResponse(&w, err, http.StatusBadRequest) != nil {
		return
	}
	err = reqData.Decode(body)
	if errorResponse(&w, err, http.StatusBadRequest) != nil {
		return
	}

	key := reqData.Key
	server.mutex.RLock()
	value, ok := server.inMemoryStorage.Kv[string(key)]
	server.mutex.RUnlock()
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

// This function must be called only when
// srv.mutex is locked
func persistUpdateWithLockedMutex(srv *serverState) error {
	srv.storageFD.Truncate(0)
	srv.storageFD.Seek(0, io.SeekStart)
	encoder := gob.NewEncoder(srv.storageFD)
	err := encoder.Encode(&srv.inMemoryStorage)
	return err
}

func sendUpdateToSynchronousFollower(w http.ResponseWriter, reqData protocol.SetRequest) error {
	//send updates to the synchronous follower
	reqBytes := reqData.Encode()
	req, err := http.NewRequest(http.MethodPut, "http://"+SYNC_FOLLOWER_URL+"/put",
		bytes.NewReader(reqBytes))
	if errorResponse(&w, err, http.StatusInternalServerError) != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := http.DefaultClient.Do(req)
	if errorResponse(&w, err, http.StatusInternalServerError) != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("got %d(%s) from synchornous replica. Not proceeding with the write",
			resp.StatusCode,
			http.StatusText(resp.StatusCode))
		errorResponse(&w, err, http.StatusInternalServerError)
		return err
	}
	return nil
}

func (server *serverState) putHandler(w http.ResponseWriter, r *http.Request) {
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
	server.mutex.Lock()
	_, keyExists := server.inMemoryStorage.Kv[string(reqData.Key)]
	if server.mode == PRIMARY {
		if addToWAL(server, string(reqData.Key), string(reqData.Value), w) != nil {
			return
		}

		if sendUpdateToSynchronousFollower(w, reqData) != nil {
			return
		}
	}
	server.inMemoryStorage.Kv[string(reqData.Key)] = string(reqData.Value)
	err = persistUpdateWithLockedMutex(server)
	server.mutex.Unlock()
	if errorResponse(&w, err, http.StatusInternalServerError) != nil {
		return
	}
	if keyExists {
		w.Write([]byte("UPDATED"))
	} else {
		w.Write([]byte("INSERTED"))
	}
}

func (srv *serverState) readStorage(filename string) {
	f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error reading storage %s", filename))
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&srv.inMemoryStorage)
	if err != io.EOF && err != nil {
		panic(fmt.Sprintf("Error decoding storage %s: %s", filename, err))
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: server <primary|primary_partition node_number|sync_follower|async_follower>\n")
	os.Exit(1)
}

func newServer(mode serverMode, url string) *serverState {
	var walFD *os.File
	var err error
	var server *serverState
	if mode == PRIMARY {
		walFD, err = os.OpenFile(WAL_FILEPATH, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
		if err != nil {
			panic(fmt.Sprintf("Error opening WAL file %s", WAL_FILEPATH))
		}
	}
	storageFilename := STORAGE_FILE_PREFIX + fmt.Sprint(mode)
	server = &serverState{
		walFD:         walFD,
		walGobEncoder: gob.NewEncoder(walFD),
		mode:          mode,
		inMemoryStorage: inMemoryStorage{
			TransactionID: 0,
			Kv:            map[string]string{},
		},
		storageFD: nil,
	}
	server.readStorage(storageFilename)
	server.storageFD, err = os.OpenFile(storageFilename, os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		panic(err)
	}
	return server
}

func (srv *serverState) storageFilename() string {
	return STORAGE_FILE_PREFIX + fmt.Sprint(srv.mode)

}

func (srv *serverState) shutDown() {
	if srv.mode == PRIMARY {
		srv.walFD.Close()
	}
	srv.storageFD.Close()
}

func main() {
	nArgs := len(os.Args)
	if nArgs < 2 {
		usage()
	}
	var url string

	var err error
	var mode serverMode
	var nodeNumber int
	switch os.Args[1] {
	case "primary":
		mode = PRIMARY
		url = SERVER + ":" + fmt.Sprint(PRIMARY_PORT)
	case "primary_partition":
		if nArgs != 3 {
			usage()
		}
		nodeNumber, err = strconv.Atoi(os.Args[2])
		if err != nil {
			usage()
		}
		mode = PRIMARY_SINGLE
		url = SERVER + ":" + fmt.Sprint(PRIMARY_PORT+PARTITION_PORT_OFFSET*nodeNumber)

	case "sync_follower":
		mode = SYNCHRONOUS_FOLLOWER
		url = SYNC_FOLLOWER_URL
	case "async_follower":
		mode = ASYNCHRONOUS_FOLLOWER
		url = ASYNC_FOLLOWER_URL
	default:
		usage()
	}

	server := newServer(mode, url)
	defer server.shutDown()

	fmt.Printf("Welcome to the distributed K-V store server in %s mode\n", server.mode)
	fmt.Printf("Listening on: %s\n", url)
	fmt.Printf("Using storage file: %s\n", server.storageFilename())

	http.HandleFunc("/get", server.getHandler)
	// Idea is that writes should be accepted only by the primary node,
	// but it is not enforced, since we do not have any auth anyway
	if mode != ASYNCHRONOUS_FOLLOWER {
		http.HandleFunc("/put", server.putHandler)
	}

	if mode == ASYNCHRONOUS_FOLLOWER {
		go server.pullPrimaryForUpdatesForever()
	}

	http.HandleFunc("/async-catchup", server.walGetHandler)

	http.ListenAndServe(url, nil)
}

// TODO: add partitioning function with hardcoded number of nodes
// TODO: make server to check the partition number and forward requests
