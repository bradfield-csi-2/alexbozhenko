package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	listeningSocketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	defer unix.Close(listeningSocketFD)
	handleError(err)
	err = unix.Bind(listeningSocketFD, &unix.SockaddrInet4{
		Port: 1025,
		Addr: [4]byte{0, 0, 0, 0},
	})
	handleError(err)
	err = unix.Listen(listeningSocketFD, 0)
	handleError(err)
	connectionSocketFD, clientSockAddr, err := unix.Accept(listeningSocketFD)
	handleError(err)
	fmt.Println(connectionSocketFD, clientSockAddr)
	defer unix.Close(connectionSocketFD)

	// same size as sysctl net.core.rmem_max
	// buf := make([]byte, 212992)
	buf := make([]byte, 30)

	for {
		//3time.Sleep(190 * time.Second)
		nBytesRead, _, err := unix.Recvfrom(connectionSocketFD, buf, 0)
		//fmt.Printf("len=%d cap=%d %v\n", len(buf), cap(buf), buf)
		//fmt.Println(hex.Dump(buf))
		fmt.Println(string(buf))
		handleError(err)
		if nBytesRead == 0 {
			// according to man 2 recv,
			// When a stream socket peer has performed an orderly shutdown, the return value will be 0 (the traditional "end-of-
			// file" return).
			break
		}
		err = unix.Sendto(connectionSocketFD, buf[:nBytesRead], 0, clientSockAddr)
		handleError(err)
		//fmt.Printf("len=%d cap=%d %v\n", len(buf), cap(buf), buf)
	}

}
