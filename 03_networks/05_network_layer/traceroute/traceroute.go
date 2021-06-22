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
	"unsafe"
)

type icmp_header struct {
	Icmp_type byte   // 8 for echo message;	// 0 for echo reply message. //11 TTL exceeded
	Code      byte   // 0
	Checksum  uint16 // The checksum is the 16-bit ones's complement of the one's
	// complement sum of the ICMP message starting with the ICMP Type.
	// For computing the checksum , the checksum field should be zero.
	// If the total length is odd, the received data is padded with one
	// octet of zeros for computing the checksum.  This checksum may be
	// replaced in the future.
}

type icmp_echo struct {
	Identifier uint16 // If code = 0, an identifier to aid in matching echos and replies,
	// may be zero.
	Seq uint16
	// 	  If code = 0, a sequence number to aid in matching echos and
	// 	  replies, may be zero.
	Payload [56]byte // Picked the size from 'man ping'
}

// Wait for ICMP response with given id and seq numbers,
// until timeout channel is populated
func receiveICMP(identifier uint16, seq uint16, timeout chan bool,
	gotResponseFromDestanationHost chan bool) {
	netaddr, _ := net.ResolveIPAddr("ip4", "0.0.0.0")
	conn, _ := net.ListenIP("ip4:icmp", netaddr)
	defer conn.Close()

	for {
		select {
		case <-timeout:
			fmt.Printf("Timeout waiting for matching ICMP response with seq=%v\n", seq)
			return
		default:
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			recieveBuf := make([]byte, 4096)
			oobBuf := make([]byte, 4096)
			numRead, oobNumRead, _, srcAddr, _ := conn.ReadMsgIP(recieveBuf, oobBuf)

			versionAndLength := recieveBuf[0]
			ipHeaderLengthBytes := 4 * (byte(versionAndLength) & byte(0xf))
			recieveBuf = recieveBuf[ipHeaderLengthBytes:]

			numRead -= int(ipHeaderLengthBytes)
			if numRead > 0 || oobNumRead > 0 {
				r := bytes.NewReader(recieveBuf[:numRead])

				var icmp_header icmp_header

				if err := binary.Read(r, binary.BigEndian, &icmp_header); err != nil {
					fmt.Println("binary.Read failed!", err)
				}
				recieveBuf = recieveBuf[unsafe.Sizeof(icmp_header):]

				if icmp_header.Icmp_type == 11 { //TTL Exceeded

					unused_bytes := 4
					recieveBuf = recieveBuf[unused_bytes:]
					recieveBuf = recieveBuf[ipHeaderLengthBytes:]

					recieveBuf = recieveBuf[unsafe.Sizeof(icmp_header):]
					var time_exceeded_msg_original_icmp_echo icmp_echo
					ttl_exceeded_reader := bytes.NewReader(recieveBuf)

					if err := binary.Read(ttl_exceeded_reader, binary.BigEndian, &time_exceeded_msg_original_icmp_echo); err != nil {
						fmt.Println("binary.Read failed!", err)
					}

					// if id and seq match, then we are done with this reply
					if time_exceeded_msg_original_icmp_echo.Identifier == identifier && time_exceeded_msg_original_icmp_echo.Seq == seq {
						fmt.Println("srcaddr from router", srcAddr)
						return
					} else {
						fmt.Printf("got packet with id=%v seq=%v, but was expecting is=%v seq=%v\n", time_exceeded_msg_original_icmp_echo.Identifier,
							time_exceeded_msg_original_icmp_echo.Seq,
							identifier,
							seq,
						)
					}

				} else { // assuming echo repsly, LOL)
					fmt.Println("srcaddr from host!", srcAddr)
					gotResponseFromDestanationHost <- true
					return
				}

			} else {
				fmt.Println("Timeout waiting for ICMP response(5s), will retry")
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
	defer ipconn.Close()

	identifier := uint16(rand.Intn(0xffff))
	seq := uint16(rand.Intn(0xffff))

	gotResponseFromDestanationHost := make(chan bool)

	for TTL := 1; TTL <= 30; TTL++ {
		select {
		case <-gotResponseFromDestanationHost:
			return
		default:
			icmp_header := &icmp_header{
				Icmp_type: 8, //echo request
				Code:      0,
				Checksum:  0,
			}
			icmp_echo := &icmp_echo{
				Identifier: identifier,
				Seq:        seq,
			}
			copy(icmp_echo.Payload[:], "Hello, 🐈!")
			buf := &bytes.Buffer{}
			binary.Write(buf, binary.BigEndian, icmp_header)
			binary.Write(buf, binary.BigEndian, icmp_echo)
			icmp_header.Checksum = checksum(buf.Bytes())
			buf = &bytes.Buffer{}
			binary.Write(buf, binary.BigEndian, icmp_header)
			binary.Write(buf, binary.BigEndian, icmp_echo)

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

			fmt.Println("===========================")
			fmt.Printf("Sent packet id=%v, seq=%v, ttl=%v\n", identifier, seq, TTL)
			timeout := make(chan bool, 1)
			go receiveICMP(identifier, seq, timeout, gotResponseFromDestanationHost)
			time.Sleep(10 * time.Second)
			timeout <- true
			seq++
		}
	}
}
