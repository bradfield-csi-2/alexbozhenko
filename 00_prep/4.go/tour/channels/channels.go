package main

import (
	"fmt"
	"time"
)

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}

	c <- sum // send sum to c
	now := time.Now()
	fmt.Println("array:", s, "the current datetime is:", now)
}

func main() {
	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x := <-c
	time.Sleep(1 * time.Second)
	y := <-c // receive from c

	time.Sleep(3 * time.Second)
	fmt.Println(x, y, x+y)
}
