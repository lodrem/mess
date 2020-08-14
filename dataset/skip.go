package dataset

import "math/rand"

func Skip(rate int) bool {
	return rand.Intn(100) < rate
}
