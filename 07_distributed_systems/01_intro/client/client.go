package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

const (
	SERVER = "127.0.0.1"
	PORT   = "8000"
	URL    = "http://" + SERVER + ":" + PORT
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func get(key string) (response string, err error) {
	return "result of get " + key, nil
	//TODO: re-use tcp connection for entire repl session?
	resp, err := http.Get(URL)
	must(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func set(key, value string) (response string, err error) {
	return "result of set " + key + ":" + value, nil
	jsonData, err := json.Marshal(struct {
		key   string
		value string
	}{key: key,
		value: value})

	// TODO: bytes.NewReader copies the []byte, right?
	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	must(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), err
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

	fmt.Println("Welcome to distributed K-V store")
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
