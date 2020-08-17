package dataset

import (
	"log"
	"math/big"
	"math/rand"
)

func IntRange(min, max *big.Int) *big.Int {

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("failed to generate data for IntRange(%d, %d): %s", min, max, r)
		}
	}()

	d := new(big.Int)
	d.Sub(max, min)
	d.Div(d, big.NewInt(2))



	r := new(big.Int).Set(min)
	if d.Int64() == 0 {
		return r
	}
	r.Add(r, big.NewInt(rand.Int63n(d.Int64())))
	r.Add(r, big.NewInt(rand.Int63n(d.Int64())))

	return r
}
