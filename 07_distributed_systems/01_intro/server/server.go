package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"kv_store/helpers"
	"kv_store/protocol"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/clientv3"
)

const (
	SERVER                     = "127.0.0.1"
	PRIMARY_PORT               = 8000
	PARTITION_PORT_OFFSET      = 1000
	SYNCHRONOUS_FOLLOWER_PORT  = "8001"
	ASYNCHRONOUS_FOLLOWER_PORT = "8002"
	SYNC_FOLLOWER_URL          = "http://" + SERVER + ":" + SYNCHRONOUS_FOLLOWER_PORT
	ASYNC_FOLLOWER_URL         = "http://" + SERVER + ":" + ASYNCHRONOUS_FOLLOWER_PORT
	STORAGE_FILE_PREFIX        = "storage"
	WAL_FILEPATH               = "wal"

	//Consistent Hash Ring config:
	CHR_NODES_NUMBER        = 3
	CHR_PARTITIONS_PER_NODE = 100
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
	hashRing        *consistentHashRing
	url             string
	etcdClient      *clientv3.Client
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
	if helpers.ErrorResponse(&w, err, http.StatusInternalServerError) != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := http.DefaultClient.Do(req)
	if helpers.ErrorResponse(&w, err, http.StatusInternalServerError) != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("got %d(%s) from synchornous replica. Not proceeding with the write",
			resp.StatusCode,
			http.StatusText(resp.StatusCode))
		helpers.ErrorResponse(&w, err, http.StatusInternalServerError)
		return err
	}
	return nil
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

func newServer(mode serverMode, urlString string) *serverState {
	var walFD *os.File
	var err error
	var server *serverState
	if mode == PRIMARY {
		walFD, err = os.OpenFile(WAL_FILEPATH, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
		helpers.PanicOnError(err)
	}
	port := ""
	u, err := url.ParseRequestURI(urlString)
	helpers.PanicOnError(err)
	if mode == PRIMARY_PARTITION {
		port = u.Port()
	}
	storageFilename := strings.Join([]string{STORAGE_FILE_PREFIX, fmt.Sprint(mode), port}, "_")

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"0.0.0.0:2379", "0.0.0.0:22379", "0.0.0.0:32379"},
		DialTimeout: 5 * time.Second,
	})
	helpers.PanicOnError(err)
	defer etcdClient.Close()
	server = &serverState{
		walFD:         walFD,
		walGobEncoder: gob.NewEncoder(walFD),
		mode:          mode,
		inMemoryStorage: inMemoryStorage{
			TransactionID: 0,
			Kv:            map[string]string{},
		},
		storageFD:  nil,
		hashRing:   NewConsistentHashRing(CHR_PARTITIONS_PER_NODE, CHR_NODES_NUMBER),
		url:        urlString,
		etcdClient: etcdClient,
	}
	for i := 0; i < CHR_NODES_NUMBER; i++ {
		server.hashRing.addNode(
			fmt.Sprintf("node:%d", i),
			getNodeUrl(i).String(),
		)
	}
	server.readStorage(storageFilename)
	server.storageFD, err = os.OpenFile(storageFilename, os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		panic(err)
	}
	return server
}

func (srv *serverState) storageFilename() string {
	return srv.storageFD.Name()
}

func (srv *serverState) shutDown() {
	if srv.mode == PRIMARY {
		srv.walFD.Close()
	}
	srv.storageFD.Close()
}

func getNodeUrl(nodeNumber int) *url.URL {
	url, err := url.Parse("http://" + SERVER + ":" +
		fmt.Sprint(PRIMARY_PORT+PARTITION_PORT_OFFSET*nodeNumber))
	helpers.PanicOnError(err)
	return url

}

func main() {
	nArgs := len(os.Args)
	if nArgs < 2 {
		usage()
	}
	var serverUrl *url.URL

	var err error
	var mode serverMode
	var nodeNumber int
	switch os.Args[1] {
	case "primary":
		mode = PRIMARY
		serverUrl, err = url.Parse(SERVER + ":" + fmt.Sprint(PRIMARY_PORT))
		helpers.PanicOnError(err)
	case "primary_partition":
		if nArgs != 3 {
			usage()
		}
		nodeNumber, err = strconv.Atoi(os.Args[2])
		if err != nil {
			usage()
		}
		mode = PRIMARY_PARTITION
		serverUrl = getNodeUrl(nodeNumber)

	case "sync_follower":
		mode = SYNCHRONOUS_FOLLOWER
		serverUrl, err = url.Parse(SYNC_FOLLOWER_URL)
		helpers.PanicOnError(err)
	case "async_follower":
		mode = ASYNCHRONOUS_FOLLOWER
		serverUrl, err = url.Parse(ASYNC_FOLLOWER_URL)
		helpers.PanicOnError(err)
	default:
		usage()
	}

	server := newServer(mode, serverUrl.String())
	defer server.shutDown()

	fmt.Printf("Welcome to the distributed K-V store server in %s mode\n", server.mode)
	fmt.Printf("Listening on: %s\n", serverUrl)
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

	http.ListenAndServe(serverUrl.Host, nil)
}
