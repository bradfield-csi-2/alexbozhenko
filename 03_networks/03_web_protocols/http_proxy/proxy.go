package main

import (
	"fmt"
	"sync"

	"golang.org/x/sys/unix"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func handleClient(listeningSocketFD int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	fmt.Printf("waitGroup: %v\n", waitGroup)

	connectionSocketFD, clientSockAddr, err := unix.Accept(listeningSocketFD)
	handleError(err)
	defer unix.Close(connectionSocketFD)

	waitGroup.Add(1)
	go handleClient(listeningSocketFD, waitGroup)

	// same size as sysctl net.core.rmem_max
	buf := make([]byte, 212992)

	for {
		nBytesRead, _, err := unix.Recvfrom(connectionSocketFD, buf, 0)
		handleError(err)
		if nBytesRead == 0 {
			// according to man 2 recv,
			// When a stream socket peer has performed an orderly shutdown, the return value will be 0 (the traditional "end-of-
			// file" return).
			break
		}
		err = unix.Sendto(connectionSocketFD, buf[:nBytesRead], 0, clientSockAddr)
		handleError(err)
	}

}

func main() {
	var wg sync.WaitGroup
	listeningSocketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	defer unix.Close(listeningSocketFD)
	handleError(err)
	err = unix.Bind(listeningSocketFD, &unix.SockaddrInet4{
		Port: 1025,
		Addr: [4]byte{0, 0, 0, 0},
	})
	handleError(err)
	err = unix.Listen(listeningSocketFD, 1)
	handleError(err)
	wg.Add(1)
	go handleClient(listeningSocketFD, &wg)
	wg.Wait()
}
