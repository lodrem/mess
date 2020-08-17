package dataset

import (
	"log"
	"math/rand"
)

func Float(precision, scale int) float64 {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("failed to generate data for Float(%d, %d): %s", precision, scale, r)
		}
	}()
	max := 0
	for i := 0; i < precision-scale; i++ {
		max = max*10 + 9
	}
	min := - max

	if max == 0 {
		return rand.Float64()
	}

	return float64(int64(min)+rand.Int63n(int64(max-min))) + rand.Float64()
}
