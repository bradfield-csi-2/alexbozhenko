package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/sys/unix"
)

const HTTP_HEADER_DELIMETER string = "\r\n\r\n"

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func handleClient(proxyListeningSocketFD int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	fmt.Printf("waitGroup: %v\n", waitGroup)

	proxyIncomingConnectionSocketFD, clientSockAddr, err := unix.Accept(proxyListeningSocketFD)
	defer unix.Close(proxyIncomingConnectionSocketFD)
	handleError(err)

	// As soon as connection was accepted by listen(),
	// we just create another goroutine that will block on listen(),
	// waiting for next connection to accept.
	// So we have unbounded parallelism, serving as many clients as
	// operating system limits allow
	waitGroup.Add(1)
	go handleClient(proxyListeningSocketFD, waitGroup)

	destinationSockAddr := &unix.SockaddrInet4{
		Port: 9000,
		Addr: [4]byte{127, 0, 0, 1},
	}

	destinationSocketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	defer unix.Close(destinationSocketFD)
	handleError(err)
	err = unix.Connect(destinationSocketFD, destinationSockAddr)
	handleError(err)

	// same size as sysctl net.core.rmem_max
	buf := make([]byte, 212992)
	sendBuffer := []byte{}

	for {
		nBytesRead, _, err := unix.Recvfrom(proxyIncomingConnectionSocketFD, buf, 0)
		handleError(err)
		if nBytesRead == 0 {
			// according to man 2 recv,
			// When a stream socket peer has performed an orderly shutdown, the return value will be 0 (the traditional "end-of-
			// file" return).
			break
		}
		// having tcp a stream of bytes, we collect the stream to sendBuffer
		// until we found http headers delimeter
		// (assuming that we have a GET request, we have no body)
		if delimeterIndex := bytes.LastIndex(buf[:nBytesRead],
			[]byte(HTTP_HEADER_DELIMETER)); delimeterIndex == -1 {
			sendBuffer = append(sendBuffer, buf[:nBytesRead]...)
		} else {
			sendBuffer = append(sendBuffer, buf[:delimeterIndex+len(HTTP_HEADER_DELIMETER)]...)
			request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(sendBuffer)))
			handleError(err)
			fmt.Println(request)
			err = unix.Sendto(destinationSocketFD, sendBuffer, 0, clientSockAddr)
			handleError(err)
			sendBuffer = buf[delimeterIndex+len(HTTP_HEADER_DELIMETER) : nBytesRead]
		}
	}

}

func main() {
	var wg sync.WaitGroup
	proxyListeningSocketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	defer unix.Close(proxyListeningSocketFD)
	handleError(err)
	err = unix.Bind(proxyListeningSocketFD, &unix.SockaddrInet4{
		Port: 8000,
		Addr: [4]byte{0, 0, 0, 0},
	})
	handleError(err)

	// Let's close the socket using RST, so we do not have to wait for sockets in
	// TIME_WAIT to expire when the program is closed, and we want to restart
	// immediately.
	err = unix.SetsockoptLinger(proxyListeningSocketFD, unix.SOL_SOCKET, unix.SO_LINGER,
		&unix.Linger{
			Onoff:  1,
			Linger: 0,
		})
	handleError(err)

	err = unix.Listen(proxyListeningSocketFD, 1)
	handleError(err)
	wg.Add(1)
	go handleClient(proxyListeningSocketFD, &wg)
	wg.Wait()
}
