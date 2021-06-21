package main

import (
	"fmt"
	"unsafe"
)

type empty_interface = interface{}

func get_int_val(i empty_interface) byte {

	// Using type assertion:
	// asserted := i.(string)
	// fmt.Println(asserted)
	// return asserted

	// Using type switch:
	// switch t := i.(type) {
	// case int:
	// 	return t
	// default:
	// 	panic(fmt.Sprintf("unexpected type %T: %v", i, i))
	// }

	//using iface internal struct
	// https://golang.org/src/runtime/runtime2.go#L203
	type iface struct {
		tab  unsafe.Pointer
		data unsafe.Pointer
	}

	iface_struct := (*iface)(unsafe.Pointer(&i))
	// for
	//	return *(*int)(iface_struct.data)
	// If concrete type size is <= uintptr, iface_struct.data would store the value directly
	// so return would look like:
	return (*((*byte)(iface_struct.data)))
}

func main() {
	var i empty_interface
	i = byte(255)
	//i := empty_interface(9000)
	fmt.Println(get_int_val(i))

}
