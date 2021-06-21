package main

/*
#include "square.h"
*/
import "C"
import "fmt"

func main() {
	x := int(C.square(2))
	fmt.Println(x)
}
