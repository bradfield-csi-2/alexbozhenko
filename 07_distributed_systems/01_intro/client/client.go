package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"kv_store/protocol"
	"net/http"
	"os"
	"regexp"
	"time"
)

const (
	SERVER = "127.0.0.1"
	PORT   = "8000"
	URL    = "http://" + SERVER + ":" + PORT
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func get(key string) (response string, err error) {
	//TODO: re-use tcp connection for entire repl session?
	getRequest := protocol.GetRequest{
		Key: []byte(key),
	}
	reqBytes := getRequest.Encode()

	req, err := http.NewRequest(http.MethodGet, URL+"/get", bytes.NewReader(reqBytes))
	panicOnError(err)
	req.Header.Set("Content-Type", "application/octet-stream")

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	elapsed := time.Since(start)
	panicOnError(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body) + " <" + fmt.Sprintf("%4.2f", float64(elapsed)/1_000_000.0) + "ms>", err
}

func set(key, value string) (response string, err error) {
	setRequest := protocol.SetRequest{
		Key:   []byte(key),
		Value: []byte(value),
	}
	reqBytes := setRequest.Encode()
	// bytes.NewReader copies the []byte, right?
	// No, it does not! https://stackoverflow.com/a/39993797/1572363
	// You should've known.

	req, err := http.NewRequest(http.MethodPut, URL+"/put",
		bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/octet-stream")
	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	elapsed := time.Since(start)
	panicOnError(err)
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

func main() {
	getRE := regexp.MustCompile(`(get) ([^=]+)$`)
	putRE := regexp.MustCompile(`(set) ([^=]+)=(.*)`)

	fmt.Println("Welcome to the distributed K-V store client")
	fmt.Println(`We support the following syntax:
get foo
set foo=bar`)
	var i int = 1
	fmt.Printf("\nIn [%d]: ", i)

	var verb, key, value, resp string
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		verb = ""
		verb, key, value = parseCommand(getRE, putRE, scanner.Bytes())
		fmt.Printf("Out[%d]: ", i)

		switch {
		case verb == "get":
			resp, err = get(key)
		case verb == "set":
			resp, err = set(key, value)
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
