package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

//https://datatracker.ietf.org/doc/html/rfc1035#section-4.1
// The header contains the following fields:
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      ID                       |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    QDCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ANCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    NSCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ARCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type dnsHeader struct {
	id      uint16
	flags   uint16
	qdcount uint16
	ancount uint16
	nscount uint16
	arcount uint16
}

type question struct {
	qname  []byte
	qtype  uint16 // https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.2
	qclass uint16 // https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.4
}

func hostnameToQname(hostname string) []byte {
	qname := []byte{}

	for _, label := range strings.Split(hostname, ".") {
		qname = append(qname, byte(len(label)))
		qname = append(qname, []byte(label)...)
	}
	qname = append(qname, byte(0))
	return qname
}

func main() {
	cloudFareDNSSockAddr := unix.SockaddrInet4{
		Port: 53,
		Addr: [4]byte{1, 1, 1, 1}, //cloudfare DNS
	}

	args := os.Args[1:]
	var hostname string

	if len(args) != 1 {
		panic("Usage: dns_client HOSTNAME")

	} else {
		hostname = os.Args[1]
	}

	socketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		panic(err)
	}

	err = unix.Sendto(socketFD, []byte("asdfasd"), 0, &cloudFareDNSSockAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println(err)

	fmt.Println(hostname)

	time.Sleep(30 * time.Second)
}
