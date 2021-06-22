package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// Internent checksum
//https://datatracker.ietf.org/doc/html/rfc1071.html
func checksum(buf []byte) uint16 {
	sum := uint32(0)

	// (1)  Adjacent octets to be checksummed are paired to form 16-bit
	// integers, and the 1's complement sum of these 16-bit integers is
	// formed.
	for ; len(buf) >= 2; buf = buf[2:] {
		sum += uint32(buf[0])<<8 | uint32(buf[1])
	}

	// If the total length is odd, the received data is padded with one
	// octet of zeros for computing the checksum.  This checksum may be
	// replaced in the future.
	if len(buf) > 0 {
		sum += uint32(buf[0]) << 8
	}

	// On a 2's complement machine, the 1's complement sum must be
	// computed by means of an "end around carry", i.e., any overflows
	// from the most significant bits are added into the least
	// significant bits. See the examples below.
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}

	csum := ^uint16(sum)
	return csum
}

func main() {
	args := os.Args[1:]
	var host string

	if len(args) != 1 {
		panic("Usage: traceroute HOST")

	} else {
		host = os.Args[1]
	}

	fmt.Println("hello", host)
	ipaddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		panic(err)
	}
	fmt.Println(ipaddr)
	ipconn, err := net.DialIP("ip:icmp", nil, ipaddr)
	if err != nil {
		panic(err)
	}
	fmt.Println(ipconn)

	type icmp_echo struct {
		icmp_type byte   // 8 for echo message;	// 0 for echo reply message.
		code      byte   // 0
		checksum  uint16 // The checksum is the 16-bit ones's complement of the one's
		// complement sum of the ICMP message starting with the ICMP Type.
		// For computing the checksum , the checksum field should be zero.
		// If the total length is odd, the received data is padded with one
		// octet of zeros for computing the checksum.  This checksum may be
		// replaced in the future.
		identifier uint16 // If code = 0, an identifier to aid in matching echos and replies,
		// may be zero.
		seq uint16
		// 	  If code = 0, a sequence number to aid in matching echos and
		// 	  replies, may be zero.
		data [56]byte
	}

	echo_msg := &icmp_echo{
		icmp_type:  8,
		code:       0,
		checksum:   0,
		identifier: 0,
		seq:        0,
	}
	copy(echo_msg.data[:], "Hello, 🐈!")
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, echo_msg)
	echo_msg.checksum = checksum(buf.Bytes())
	buf = &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, echo_msg)
	fmt.Println(buf.Len())
	fmt.Println(buf.Bytes())

	// TODO: increase ttl starting from 1

	// TODO: parse responses

	ipconn.Write(buf.Bytes())

}
