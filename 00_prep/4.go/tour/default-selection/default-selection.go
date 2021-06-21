package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	t := time.Tick(time.Second)
	fmt.Println(reflect.TypeOf(t))
	fmt.Println(<-t)
	fmt.Println(<-t)
}
