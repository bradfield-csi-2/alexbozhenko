package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
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

// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.3
// Each resource record has the following format:
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                                               |
//     /                                               /
//     /                      NAME                     /
//     |                                               |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      TYPE                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     CLASS                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      TTL                      |
//     |                                               |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                   RDLENGTH                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
//     /                     RDATA                     /
//     /                                               /
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type resourceRecordHeader struct {
	Name     uint16
	Type     uint16
	Class    uint16
	TTL      uint32
	RdLength uint16
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

func generateQueryHeader(id uint16) dnsHeader {
	header := dnsHeader{
		ID:                 id,
		QR_Opcode_AA_TC_RD: 1, // rd = 1
		RA_Z_Rcode:         0,
		Qdcount:            uint16(1), //We always send one query
		Ancount:            0,
		Nscount:            0,
		Arcount:            0,
	}
	return header
}

func generateQueryMessage(hostname string) question {
	return question{
		qname:  hostnameToQname(hostname),
		qtype:  TYPE_A,
		qclass: CLASS_IN,
	}
}

func generateQueryPayload(hostname string) ([]byte, uint16, int) {
	buffer := &bytes.Buffer{}
	id := uint16(rand.Uint32())
	header := generateQueryHeader(id)
	binary.Write(buffer, binary.BigEndian, header)
	queryMessage := generateQueryMessage(hostname)

	// Why does go-staticcheck says:
	// "value cannot be used with binary.Write (SA1003)"
	// when I try to write entire struct, like in the commented line below
	// I guess that's because byte slice is stored as a pointer in the struct?
	// Is writing each field individually is the right way to do that?
	//binary.Write(buffer, binary.BigEndian, q)
	binary.Write(buffer, binary.BigEndian, queryMessage.qname)
	binary.Write(buffer, binary.BigEndian, queryMessage.qtype)
	binary.Write(buffer, binary.BigEndian, queryMessage.qclass)

	queryMessageLength := (binary.Size(queryMessage.qname) +
		binary.Size(queryMessage.qtype) +
		binary.Size(queryMessage.qclass))

	return buffer.Bytes(), id, queryMessageLength
}

func main() {
	cloudFareDNSSockAddr := unix.SockaddrInet4{
		Port: 53,
		Addr: [4]byte{1, 1, 1, 1}, //cloudfare DNS
	}

	args := os.Args[1:]
	var hostname string

	if len(args) != 1 {
		// actually, multiple queries in the same message are not supported
		// https://stackoverflow.com/a/4083071/1572363
		panic("Usage: dns_client HOSTNAME")
	} else {
		hostname = os.Args[1]
	}

	socketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		panic(err)
	}

	query, requestId, queryMessageLength := generateQueryPayload(hostname)

	err = unix.Sendto(socketFD, query, 0, &cloudFareDNSSockAddr)
	if err != nil {
		panic(err)
	}

	// TODO: In real word, what should be the buffer size, having that
	// we do not know the size of the reponse beforehand?
	// For now, setting it to max possible size of the udp message
	dnsResponse := make([]byte, 65535)
	nBytesRead, _, err := unix.Recvfrom(socketFD, dnsResponse, 0)
	if err != nil {
		panic(err)
	}
	dnsResponse = dnsResponse[:nBytesRead]

	responseHeader := dnsHeader{}
	buf := bytes.NewReader(dnsResponse)
	err = binary.Read(buf, binary.BigEndian, &responseHeader)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	if requestId != responseHeader.ID {
		panic("Wrong request id in the response that we got. Ignoring!")
	}
	// ignoring original request at the beginning of response
	buf.Seek(int64(queryMessageLength), io.SeekCurrent)

	for rrNumber := uint16(0); rrNumber < responseHeader.Ancount; rrNumber++ {
		rrHeader := resourceRecordHeader{}
		err = binary.Read(buf, binary.BigEndian, &rrHeader)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}

		ip := make([]byte, rrHeader.RdLength)
		err = binary.Read(buf, binary.BigEndian, &ip)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		fmt.Println(net.IPv4(ip[0], ip[1], ip[2], ip[3]))

	}
}
