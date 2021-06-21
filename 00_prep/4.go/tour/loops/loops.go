package main

import (
	"fmt"
	"math"
)

// Sqrt meh
func Sqrt(x float64) float64 {
	z := 1.0
	for i := 0; i < 10; i++ {

		newGuess := z - (z*z-x)/(2*z)
		fmt.Printf("i=%v z=%v new_guess=%v\n", i, z, newGuess)

		if math.Abs((z - newGuess)) < 0.001 {
			return newGuess

		}
		z = newGuess
	}
	return z
}

func main() {
	fmt.Println(Sqrt(77))
}
