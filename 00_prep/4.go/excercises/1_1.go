package main

import (
"fmt"
"os"
//"strings"
)


func main() {
	for index, cli_arg := range os.Args {
		fmt.Printf("%d=%s\n", index, cli_arg)
	}
//fmt.Println(strings.Join(os.Args[0:], " "))
}

