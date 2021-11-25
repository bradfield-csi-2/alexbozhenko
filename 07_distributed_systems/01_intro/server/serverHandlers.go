package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"kv_store/helpers"
	"kv_store/protocol"
	"log"
	"net/http"
	"runtime"
)

func redirect(url, method string, body []byte, w *http.ResponseWriter) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if helpers.ErrorResponse(w, err, http.StatusInternalServerError) != nil {
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := http.DefaultClient.Do(req)
	if helpers.ErrorResponse(w, err, http.StatusInternalServerError) != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("got %d(%s) from redirected request",
			resp.StatusCode,
			http.StatusText(resp.StatusCode))
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	redirectedResponseBody, err := io.ReadAll(resp.Body)
	if helpers.ErrorResponse(w, err, http.StatusInternalServerError) != nil {
		return
	}
	fmt.Fprintf(*w, "%s", redirectedResponseBody)
}

func (server *serverState) getHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getHandler. Goroutines:%v", runtime.NumGoroutine())
	reqData := protocol.GetRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if helpers.ErrorResponse(&w, err, http.StatusBadRequest) != nil {
		return
	}
	err = reqData.Decode(body)
	if helpers.ErrorResponse(&w, err, http.StatusBadRequest) != nil {
		return
	}

	key := reqData.Key
	nodeNameForCurrentKey, urlForCurrentKey := server.hashRing.GetNode(string(key))
	if urlForCurrentKey != server.url {
		log.Printf("redirecting read request. key:%s to_node:%s\n",
			key, nodeNameForCurrentKey)
		redirect(urlForCurrentKey+"/get", http.MethodGet, body, &w)
		return
	}

	server.mutex.RLock()
	value, ok := server.inMemoryStorage.Kv[string(key)]
	server.mutex.RUnlock()
	if ok {
		fmt.Fprintf(w, "%s", value)
	} else {
		fmt.Fprintf(w, "key %s NOT FOUND", key)
	}
}

func (server *serverState) putHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("putHandler. Goroutines:%v", runtime.NumGoroutine())
	reqData := protocol.SetRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if helpers.ErrorResponse(&w, err, http.StatusBadRequest) != nil {
		return
	}
	err = reqData.Decode(body)
	if helpers.ErrorResponse(&w, err, http.StatusBadRequest) != nil {
		return
	}

	key := reqData.Key
	nodeNameForCurrentKey, urlForCurrentKey := server.hashRing.GetNode(string(key))
	if urlForCurrentKey != server.url {
		log.Printf("redirecting write request. key:%s to_node:%s\n",
			key, nodeNameForCurrentKey)
		redirect(urlForCurrentKey+"/put", http.MethodPut, body, &w)
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
	if helpers.ErrorResponse(&w, err, http.StatusInternalServerError) != nil {
		return
	}
	if keyExists {
		w.Write([]byte("UPDATED"))
	} else {
		w.Write([]byte("INSERTED"))
	}
}
