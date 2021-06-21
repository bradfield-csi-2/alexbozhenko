package main

import "fmt"

// ErrNegativeSqrt meh
type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

// Sqrt Newton method
func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	z := 1.0
	for i := 0; i < 10; i++ {

		newGuess := z - (z*z-x)/(2*z)

		if (newGuess - z) < 0.001 {
			return newGuess, nil
		}
		z = newGuess
	}
	return z, nil
}

func main() {
	//fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
