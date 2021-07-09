package main

import (
	"bufio"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"golang.org/x/sys/unix"
)

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
	fmt.Printf("waitGroup: %v\n", waitGroup)

	client_proxySocketFD, proxy_clientSockAddr, err :=
		unix.Accept(proxyListeningSocketFD)
	defer unix.Close(client_proxySocketFD)
	handleErrorExitGoroutine(err)
	fmt.Println(proxy_clientSockAddr)

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

	request, err := http.ReadRequest(bufio.NewReader(&client_proxySocketReaderWriter))
	handleErrorExitGoroutine(err)

	webServerSockAddr := &unix.SockaddrInet4{
		Port: 9000,
		Addr: [4]byte{127, 0, 0, 1},
	}
	webServerSocketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	defer unix.Close(webServerSocketFD)
	handleErrorExitGoroutine(err)
	err = unix.Connect(webServerSocketFD, webServerSockAddr)
	handleErrorExitGoroutine(err)
	proxy_serverSocketReaderWriter := socketReaderWriter{
		socketFD:       webServerSocketFD,
		sendToSockAddr: webServerSockAddr,
	}
	request.Write(&proxy_serverSocketReaderWriter)
	//fmt.Printf("request: %+v\n", request)
	//fmt.Println("host", request.Host)
	//fmt.Println(request.URL)

	response, err := http.ReadResponse(bufio.NewReader(&proxy_serverSocketReaderWriter), nil)
	handleErrorExitGoroutine(err)
	// stracing the program shows that body is in fact(as documented)
	// is read on demand. There is no recvfrom() syscalls for getting the body
	// unless we ask for it with io.Readall or response.Write

	// body, err := io.ReadAll(response.Body)
	// handleError(err)
	// fmt.Println(response.Header)
	// fmt.Printf("%s\n", body)

	// we can create a new buffer
	// var b bytes.Buffer
	// response.Write(&b)

	// and parse it back to then send.
	// Would it be possible to store the reponse(with body) in the map somehow,
	// instead of storing the bytes?
	// newresp, _ := http.ReadResponse(bufio.NewReader(&b), nil)
	// newresp.Write(&client_proxySocketReaderWriter)

	// Body will be closed by write
	response.Write(&client_proxySocketReaderWriter)

	// for {
	// 	nBytesRead, _, err := unix.Recvfrom(proxyIncomingConnectionSocketFD, buf, 0)
	// 	handleError(err)
	// 	if nBytesRead == 0 {
	// 		// according to man 2 recv,
	// 		// When a stream socket peer has performed an orderly shutdown, the return value will be 0 (the traditional "end-of-
	// 		// file" return).
	// 		break
	// 	}
	// 	// having tcp a stream of bytes, we collect the stream to sendBuffer
	// 	// until we found http headers delimeter
	// 	// (assuming that we have a GET request, we have no body)
	// 	if delimeterIndex := bytes.LastIndex(buf[:nBytesRead],
	// 		[]byte(HTTP_HEADER_DELIMETER)); delimeterIndex == -1 {
	// 		sendBuffer = append(sendBuffer, buf[:nBytesRead]...)
	// 	} else {
	// 		sendBuffer = append(sendBuffer, buf[:delimeterIndex+len(HTTP_HEADER_DELIMETER)]...)
	// 		request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(sendBuffer)))
	// 		handleError(err)
	// 		fmt.Println(request)
	// 		err = unix.Sendto(destinationSocketFD, sendBuffer, 0, clientSockAddr)
	// 		handleError(err)
	// 		sendBuffer = buf[delimeterIndex+len(HTTP_HEADER_DELIMETER) : nBytesRead]
	// 	}
	//}

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

	// Let's close the socket using RST, so we do not have to wait for sockets in
	// TIME_WAIT to expire when the program is closed, and we want to restart
	// immediately.
	err = unix.SetsockoptLinger(proxyListeningSocketFD, unix.SOL_SOCKET, unix.SO_LINGER,
		&unix.Linger{
			Onoff:  1,
			Linger: 0,
		})
	handleErrorPanic(err)

	err = unix.Listen(proxyListeningSocketFD, 1)
	handleErrorPanic(err)
	wg.Add(1)
	go handleClient(proxyListeningSocketFD, &wg)
	wg.Wait()
}
