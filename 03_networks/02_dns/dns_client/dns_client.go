package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

const TYPE_A = 0x0001
const CLASS_IN = 0x0001

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
	ID                 uint16
	QR_Opcode_AA_TC_RD uint8
	RA_Z_Rcode         uint8
	Qdcount            uint16
	Ancount            uint16
	Nscount            uint16
	Arcount            uint16
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
		ID:                 uint16(rand.Uint32()),
		QR_Opcode_AA_TC_RD: 1, // rd = 1
		RA_Z_Rcode:         0,
		Qdcount:            uint16(len(hostnames)),
		Ancount:            0,
		Nscount:            0,
		Arcount:            0,
	}
	buffer := &bytes.Buffer{}
	binary.Write(buffer, binary.BigEndian, header)
	var q question
	for _, hostname := range hostnames {
		q = question{
			qname:  hostnameToQname(hostname),
			qtype:  TYPE_A,
			qclass: CLASS_IN,
		}

		// Why does go-staticcheck says:
		// "value cannot be used with binary.Write (SA1003)"
		// when I try to write entire struct, like in the commented line below
		// Is writing each field individually is the right way to do that?
		//binary.Write(buffer, binary.BigEndian, q)
		binary.Write(buffer, binary.BigEndian, q.qname)
		binary.Write(buffer, binary.BigEndian, q.qtype)
		binary.Write(buffer, binary.BigEndian, q.qclass)
	}

	return buffer.Bytes()
}

func main() {
	cloudFareDNSSockAddr := unix.SockaddrInet4{
		Port: 53,
		Addr: [4]byte{1, 1, 1, 1}, //cloudfare DNS
	}

	args := os.Args[1:]
	var hostnames []string

	if len(args) <= 0 {
		// actually, multiple queries in the same message are not supported
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

	// TODO: In real word, what should be the buffer size, having that
	// we do not know the size of the reponse beforehand
	// for now, setting it to max possible size of the udp message
	dnsResponse := make([]byte, 65535)
	nBytesRead, _, err := unix.Recvfrom(socketFD, dnsResponse, 0)
	if err != nil {
		panic(err)
	}
	dnsResponse = dnsResponse[:nBytesRead]

	fmt.Println(hex.Dump(dnsResponse))
	header := dnsHeader{}
	buf := bytes.NewReader(dnsResponse)
	err = binary.Read(buf, binary.BigEndian, &header)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Printf("%+v\n", header)
	fmt.Printf("%+v\n", header)

}
