package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"kv_store/consistenHash"
	"kv_store/helpers"
	"kv_store/protocol"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ETCD_IP = "192.168.1.3"
)

func updateConsistenHashRingWithAliveNodes() {
	etcdResponse, err := etcdClient.Get(context.TODO(), "alive_servers",
		clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	for _, ev := range etcdResponse.Kvs {
		parsedUrl, err := url.Parse(string(ev.Value))
		if err == nil {
			consistenHashRing.AddNode(parsedUrl.String(), parsedUrl.String())
			fmt.Printf("Got node from etcd list of alive nodes. Adding it: %s\n", string(ev.Value))
		}
	}

}

func get(key string, consistenHashRing *consistenHash.ConsistentHashRing) (response string, err error) {
	_, url := consistenHashRing.GetNode(key)
	log.Println("Chosen url from consistent hash ring for the key:", url)

	//TODO: re-use tcp connection for entire repl session?
	getRequest := protocol.GetRequest{
		Key: []byte(key),
	}
	reqBytes := getRequest.Encode()

	req, err := http.NewRequest(http.MethodGet, url+"/get",
		bytes.NewReader(reqBytes))
	helpers.PanicOnError(err)
	req.Header.Set("Content-Type", "application/octet-stream")

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	elapsed := time.Since(start)
	helpers.PanicOnError(err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body) + " <" + fmt.Sprintf("%4.2f", float64(elapsed)/1_000_000.0) + "ms>", err
}

func set(key, value string, consistenHashRing *consistenHash.ConsistentHashRing) (response string, err error) {
	_, url := consistenHashRing.GetNode(key)
	log.Println("Chosen url from consistent hash ring for the key:", url)
	setRequest := protocol.SetRequest{
		Key:   []byte(key),
		Value: []byte(value),
	}
	reqBytes := setRequest.Encode()
	req, err := http.NewRequest(http.MethodPut, url+"/put",
		bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/octet-stream")
	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	elapsed := time.Since(start)
	helpers.PanicOnError(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body) + " <" + fmt.Sprintf("%4.2f", float64(elapsed)/1_000_000.0) + "ms>", err
}

func parseCommand(getRE, putRE *regexp.Regexp, userInput []byte) (verb, key, value string) {
	putTokens := putRE.FindSubmatch(userInput)
	getTokens := getRE.FindSubmatch(userInput)
	if len(putTokens) > 0 {
		verb = string(putTokens[1])
		key = string(putTokens[2])
		value = string(putTokens[3])
	} else if len(getTokens) > 0 {
		verb = string(getTokens[1])
		key = string(getTokens[2])
	}
	return
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: client")
	os.Exit(1)
}

func main() {
	getRE := regexp.MustCompile(`(get) ([^=]+)$`)
	putRE := regexp.MustCompile(`(set) ([^=]+)=(.*)`)

	fmt.Println("Welcome to the distributed K-V store client")
	fmt.Println(`We support the following syntax:
get foo
set foo=bar`)
	consistenHashRing := consistenHash.NewConsistentHashRing(100)
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{ETCD_IP + ":2379", ETCD_IP + ":22379", ETCD_IP + ":32379"},
		DialTimeout: 5 * time.Second,
	})
	helpers.PanicOnError(err)

	var i int = 1
	fmt.Printf("\nIn [%d]: ", i)

	var verb, key, value, resp string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		verb = ""
		verb, key, value = parseCommand(getRE, putRE, scanner.Bytes())
		fmt.Printf("Out[%d]: ", i)

		switch {
		case verb == "get":
			// TODO query etc, build hash ring
			resp, err = get(key, etcdClient, consistenHashRing)
		case verb == "set":
			// TODO query etc, build hash ring
			resp, err = set(key, value, etcdClient, consistenHashRing)
		default:
			resp = fmt.Sprintf("failed to parse command %s", scanner.Text())
		}
		fmt.Println(resp)
		if err != nil {
			fmt.Println(err)
		}

		i++
		fmt.Printf("\nIn [%d]: ", i)
	}

}
