package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"sync"

	"golang.org/x/sys/unix"
)

var KEEP_ALIVE = "Keep-Alive"

var proxyListeningSocketBacklogSize = 100000

type cachedResponse struct {
	mutex sync.RWMutex
	// Would it be possible to store the reponse(with body) in the map somehow,
	// instead of storing the bytes?
	response []byte
}

var cacheWriteMutex sync.RWMutex
var cacheMap = map[url.URL]*cachedResponse{}

var webServerSockAddr = &unix.SockaddrInet4{
	Port: 9000,
	Addr: [4]byte{127, 0, 0, 1},
}

type socketReaderWriter struct {
	socketFD       int
	sendToSockAddr unix.Sockaddr
}

func (s *socketReaderWriter) Read(p []byte) (n int, err error) {
	n, _, err = unix.Recvfrom(s.socketFD, p, 0)
	return
}

func (s *socketReaderWriter) Write(p []byte) (n int, err error) {
	err = unix.Sendto(s.socketFD, p, 0, s.sendToSockAddr)
	return len(p), err
}

func handleErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func handleErrorExitGoroutine(err error) {
	if err != nil {
		fmt.Println(err)
		runtime.Goexit()
	}
}

func handleClient(proxyListeningSocketFD int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	fmt.Printf("waitGroup currently has : %v\n", waitGroup)

	client_proxySocketFD, proxy_clientSockAddr, err :=
		unix.Accept(proxyListeningSocketFD)
	defer unix.Close(client_proxySocketFD)
	handleErrorExitGoroutine(err)

	// As soon as connection was accepted by listen(),
	// we just create another goroutine that will block on listen(),
	// waiting for next connection to accept.
	// So we have unbounded parallelism, serving as many clients as
	// operating system limits allow
	waitGroup.Add(1)
	go handleClient(proxyListeningSocketFD, waitGroup)

	client_proxySocketReaderWriter := socketReaderWriter{
		socketFD:       client_proxySocketFD,
		sendToSockAddr: proxy_clientSockAddr,
	}

	// Since we are using Keep-Alive, we are not exiting the
	// goroutine, but instead keep looping, trying to
	// read from the proxy client if klient keep sending data
	// on the same socket
	for {

		// Checking if client closed the connection
		// http://stefan.buettcher.org/cs/conn_closed.html
		// It is not clear to me how to properly check for connection
		// that is closed by client, so we can exit this goroutine.
		// This seems to work, but is it the right way?
		// Also, I had to remove MSG_DONTWAIT flag. This needs further investigation.
		bytesAvailable, _, err := unix.Recvfrom(client_proxySocketFD,
			[]byte{0}, unix.MSG_PEEK)
		handleErrorExitGoroutine(err)
		if bytesAvailable <= 0 {
			//handleErrorExitGoroutine(errors.New("closed socket detected!!!!111. Exiting this goroutine"))
			break
		}

		clientRequest, err := http.ReadRequest(bufio.NewReader(&client_proxySocketReaderWriter))
		handleErrorExitGoroutine(err)

		cacheWriteMutex.RLock()
		responseFromCache, responseFoundInCache := cacheMap[*clientRequest.URL]
		cacheWriteMutex.RUnlock()

		if responseFoundInCache {
			responseFromCache.mutex.RLock()
			parsedResponse, _ := http.ReadResponse(bufio.NewReader(
				bytes.NewReader(responseFromCache.response)), nil)

			//Set keep-alive in the response that will be sent from our proxy
			parsedResponse.Close = false
			parsedResponse.Header.Set("Connection", KEEP_ALIVE)
			parsedResponse.Write(&client_proxySocketReaderWriter) // will close the body

			//fmt.Printf("URL %v was retrieved from cache\n", clientRequest.URL.Path)
			responseFromCache.mutex.RUnlock()
		} else {
			webServerSocketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
			defer unix.Close(webServerSocketFD)
			handleErrorExitGoroutine(err)
			err = unix.Connect(webServerSocketFD, webServerSockAddr)
			handleErrorExitGoroutine(err)
			proxy_serverSocketReaderWriter := socketReaderWriter{
				socketFD:       webServerSocketFD,
				sendToSockAddr: webServerSockAddr,
			}
			clientRequest.Write(&proxy_serverSocketReaderWriter)

			serverResponse, err := http.ReadResponse(bufio.NewReader(&proxy_serverSocketReaderWriter),
				nil)
			handleErrorExitGoroutine(err)
			// strace-ing the program shows that body is in fact(as documented)
			// is read on demand. There is no recvfrom() syscalls for getting the body
			// unless we ask for it with io.Readall or response.Write

			//Set keep-alive in the response that will be sent from our proxy
			serverResponse.Close = false
			serverResponse.Header.Set("Connection", KEEP_ALIVE)

			// Write response to a buffer, so we can easly store the bite slice
			// of headers + body in proper wire format
			var responseBuffer bytes.Buffer
			serverResponse.Write(&responseBuffer) //will close response's body

			// https://golang.org/doc/faq#atomic_maps
			// Go maps can not be read while being modified
			// Locking only for period of updating the cache
			// means that parallel requests fot the same uncached URL
			// will make parallel request to the web server, but will minimise
			// time when other clients(readers) will be blocked
			cacheWriteMutex.Lock()
			cacheMap[*clientRequest.URL] = &cachedResponse{
				mutex:    sync.RWMutex{},
				response: responseBuffer.Bytes(),
			}
			cacheWriteMutex.Unlock()
			fmt.Printf("URL %v was stored in cache\n", clientRequest.URL.Path)

			// Once response was cached, return the bytes to the client
			client_proxySocketReaderWriter.Write(responseBuffer.Bytes())
		}
	}

}

func main() {
	var wg sync.WaitGroup
	proxyListeningSocketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	defer unix.Close(proxyListeningSocketFD)
	handleErrorPanic(err)
	err = unix.Bind(proxyListeningSocketFD, &unix.SockaddrInet4{
		Port: 8000,
		Addr: [4]byte{0, 0, 0, 0},
	})
	handleErrorPanic(err)

	// Let's close the connection with RST flag, so we do not have to wait for sockets in
	// TIME_WAIT to expire when the program is closed, and we want to restart
	// immediately.
	err = unix.SetsockoptLinger(proxyListeningSocketFD, unix.SOL_SOCKET, unix.SO_LINGER,
		&unix.Linger{
			Onoff:  1,
			Linger: 0,
		})
	handleErrorPanic(err)

	err = unix.Listen(proxyListeningSocketFD, proxyListeningSocketBacklogSize)
	handleErrorPanic(err)
	wg.Add(1)
	go handleClient(proxyListeningSocketFD, &wg)
	wg.Wait()
}
