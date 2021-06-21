package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

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
		icmp_type byte  // 8 for echo message;	// 0 for echo reply message.
		code      byte  // 0
		checksum  int16 // The checksum is the 16-bit ones's complement of the one's
		// complement sum of the ICMP message starting with the ICMP Type.
		// For computing the checksum , the checksum field should be zero.
		// If the total length is odd, the received data is padded with one
		// octet of zeros for computing the checksum.  This checksum may be
		// replaced in the future.
		identifier int16 // If code = 0, an identifier to aid in matching echos and replies,
		// may be zero.
		seq int16
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
	//TODO what is the correct way to cast struct to byte array
	//	echo_msg_bytes := *([]byte)(unsafe.Pointer(echo_msg))
	buf := &bytes.Buffer{}
	fmt.Println(*buf)
	binary.Write(buf, binary.LittleEndian, echo_msg)
	// TODO: calculate checksum

	// TODO: increase ttl starting from 1

	// TODO: parse responses

	ipconn.Write(buf.Bytes())
	//	fmt.Fprintf(ipconn, "GET ")

}
