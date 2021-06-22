package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"syscall"
	"time"
)

type icmp_header_fixme struct {
	Icmp_type byte   // 8 for echo message;	// 0 for echo reply message.
	Code      byte   // 0
	Checksum  uint16 // The checksum is the 16-bit ones's complement of the one's
	// complement sum of the ICMP message starting with the ICMP Type.
	// For computing the checksum , the checksum field should be zero.
	// If the total length is odd, the received data is padded with one
	// octet of zeros for computing the checksum.  This checksum may be
	// replaced in the future.
	Identifier uint16 // If code = 0, an identifier to aid in matching echos and replies,
	// may be zero.
	Seq uint16
	// 	  If code = 0, a sequence number to aid in matching echos and
	// 	  replies, may be zero.
}

type icmp_echo struct {
	Icmp_type byte   // 8 for echo message;	// 0 for echo reply message.
	Code      byte   // 0
	Checksum  uint16 // The checksum is the 16-bit ones's complement of the one's
	// complement sum of the ICMP message starting with the ICMP Type.
	// For computing the checksum , the checksum field should be zero.
	// If the total length is odd, the received data is padded with one
	// octet of zeros for computing the checksum.  This checksum may be
	// replaced in the future.
	Identifier uint16 // If code = 0, an identifier to aid in matching echos and replies,
	// may be zero.
	Seq uint16
	// 	  If code = 0, a sequence number to aid in matching echos and
	// 	  replies, may be zero.
	Payload [56]byte
}

// Wait for ICMP response with given if and seq numbers,
// until timeout channel is populated
func receiveICMP(identifier uint16, seq uint16, timeout chan bool) {
	fmt.Println("started recevie")
	netaddr, _ := net.ResolveIPAddr("ip4", "0.0.0.0")
	conn, _ := net.ListenIP("ip4:icmp", netaddr)
	defer conn.Close()

	for {
		select {
		case <-timeout:
			fmt.Println("Timeout waiting for matching ICMP response")
			return
		default:
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			recieveBuf := make([]byte, 4096)
			oobBuf := make([]byte, 4096)
			numRead, oobNumRead, _, srcAddr, _ := conn.ReadMsgIP(recieveBuf, oobBuf)

			// TODO FIXME oh... offset hardcoded ip header length
			offsetIntoOriginalData := 20 + 8 + 20
			recieveBuf = recieveBuf[offsetIntoOriginalData:]
			numRead -= offsetIntoOriginalData
			if numRead > 0 || oobNumRead > 0 {
				r := bytes.NewReader(recieveBuf[:numRead])

				var data icmp_header_fixme

				if err := binary.Read(r, binary.BigEndian, &data); err != nil {
					fmt.Println("binary.Read failed!:", err)
				}
				// fmt.Println(data.Identifier, identifier)
				// fmt.Println(data.Seq, seq)
				// fmt.Printf("%+v\n", data)
				// fmt.Printf("%x\n", recieveBuf[:numRead])
				// if id and seq match, then exit
				if data.Identifier == identifier && data.Seq == seq {
					fmt.Println("srcaddr", srcAddr)
					return
				} else {
					fmt.Printf("got packet with id=%v seq=%v\n", data.Identifier, data.Seq)
				}

			} else {
				fmt.Println("Timeout waiting for any ICMP response")
			}
		}
	}

}

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

	ipaddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		panic(err)
	}
	ipconn, err := net.DialIP("ip4:icmp", nil, ipaddr)
	if err != nil {
		panic(err)
	}

	identifier := uint16(rand.Intn(0xffff))
	seq := uint16(rand.Intn(0xffff))

	for TTL := 1; TTL <= 30; TTL++ {
		echo_msg := &icmp_echo{
			Icmp_type:  8,
			Code:       0,
			Checksum:   0,
			Identifier: identifier,
			Seq:        seq,
		}
		copy(echo_msg.Payload[:], "Hello, 🐈!")
		buf := &bytes.Buffer{}
		binary.Write(buf, binary.BigEndian, echo_msg)
		echo_msg.Checksum = checksum(buf.Bytes())
		buf = &bytes.Buffer{}
		binary.Write(buf, binary.BigEndian, echo_msg)

		fd, err := ipconn.File()
		if err != nil {
			panic(err)
		}

		// Even though https://golang.org/pkg/net/#IPConn.File says
		// {|
		// The returned os.File's file descriptor is different from the connection's.
		// Attempting to change properties of the original using this duplicate may or may not have the desired effect.
		// |}
		// This works because both fds point to the same socket, and calling
		// setsockopt on the duplicate fd propagates the change to the socket
		syscall.SetsockoptByte(int(fd.Fd()), int(syscall.IPPROTO_IP), int(syscall.IP_TTL),
			byte(TTL))
		fd.Close()
		ipconn.Write(buf.Bytes())
		fmt.Printf("Sent packet id=%v, seq=%v, ttl=%v\n", identifier, seq, TTL)
		timeout := make(chan bool, 1)
		go receiveICMP(identifier, seq, timeout)
		time.Sleep(5 * time.Second)
		timeout <- true
		seq++
	}
	ipconn.Close()

	// TODO: parse responses

}
