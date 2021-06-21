package main

import (
	"fmt"
)

type Bozhenker = interface {
	String() string
	OverNineK() string
}
type stingable_int int

func (i stingable_int) String() string {
	return fmt.Sprintf("%d\n", i)
}

func (i stingable_int) OverNineK() string {
	return fmt.Sprintf("%d\n", i+9000)
}

func show_itable(i Bozhenker) int {

	//asserted := i.(stingable_int)
	fmt.Println(i)
	fmt.Println(i.OverNineK())
	//	fmt.Println(asserted.OverNineK())
	return 0
}

func main() {
	var s Bozhenker
	s = stingable_int(42)

	fmt.Println(show_itable(s))

}
