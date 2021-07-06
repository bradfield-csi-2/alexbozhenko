package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

const host_address_type = uint16(1)
const internet_class = uint16(1)

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
	id                 uint16
	qr_Opcode_AA_TC_RD uint8
	RA_Z_Rcode         uint8
	qdcount            uint16
	ancount            uint16
	nscount            uint16
	arcount            uint16
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

func generateQueryMessage(hostnames []string) []byte {
	header := dnsHeader{
		id:                 uint16(rand.Uint32()),
		qr_Opcode_AA_TC_RD: 1, // rd = 1
		RA_Z_Rcode:         0,
		qdcount:            uint16(len(hostnames)),
		ancount:            0,
		nscount:            0,
		arcount:            0,
	}
	buffer := &bytes.Buffer{}
	binary.Write(buffer, binary.BigEndian, header)
	query := buffer.Bytes()
	b := make([]byte, 2)
	for _, hostname := range hostnames {
		query = append(query, hostnameToQname(hostname)...)
		binary.BigEndian.PutUint16(b, host_address_type)
		query = append(query, b...)
		binary.BigEndian.PutUint16(b, internet_class)
		query = append(query, b...)
	}

	return query
}

func main() {
	cloudFareDNSSockAddr := unix.SockaddrInet4{
		Port: 53,
		Addr: [4]byte{1, 1, 1, 1}, //cloudfare DNS
	}

	args := os.Args[1:]
	var hostnames []string

	if len(args) <= 0 {
		// actually multiple queries are not supported
		// https://stackoverflow.com/a/4083071/1572363
		panic("Usage: dns_client HOSTNAME [HOSTNAME]...")

	} else {
		hostnames = os.Args[1:]
	}

	socketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		panic(err)
	}

	query := generateQueryMessage(hostnames)

	err = unix.Sendto(socketFD, query, 0, &cloudFareDNSSockAddr)
	if err != nil {
		panic(err)
	}
	sockname, err := unix.Getsockname(socketFD)
	if err != nil {
		panic(err)
	}
	fmt.Println(sockname)

	dnsResponse := make([]byte, 512)
	// dnsResponse := new([]byte)

	nBytesRead, fromSock, err := unix.Recvfrom(socketFD, dnsResponse, 0)
	fmt.Println(nBytesRead, fromSock, err)
	fmt.Println(len(dnsResponse), cap(dnsResponse), dnsResponse)

	fmt.Println(hostnames)
}
