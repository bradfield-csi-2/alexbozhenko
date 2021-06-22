package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"syscall"
	"time"
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

	ipaddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		panic(err)
	}
	ipconn, err := net.DialIP("ip:icmp", nil, ipaddr)
	if err != nil {
		panic(err)
	}

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

	// TODO: increase ttl starting from 1

	fd, err := ipconn.File()
	if err != nil {
		panic(err)
	}

	ttl := byte(23)
	// Even though https://golang.org/pkg/net/#IPConn.File says
	//
	// The returned os.File's file descriptor is different from the connection's.
	// Attempting to change properties of the original using this duplicate may or may not have the desired effect.
	//
	// This works because both fds point to the same socket, and calling
	// setsockopt on the duplicate fd propagates the change to the socket
	syscall.SetsockoptByte(int(fd.Fd()), int(syscall.IPPROTO_IP), int(syscall.IP_TTL),
		ttl)
	time.Sleep(50 * time.Second)
	fd.Close()
	// error = setsockopt(send_socket, IPPROTO_IP, IP_TTL, &ttl, sizeof(ttl));

	// TODO: parse responses

	ipconn.Write(buf.Bytes())

}
