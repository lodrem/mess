package dataset

import (
	"math/big"
	"math/rand"
)

func IntRange(min, max *big.Int) *big.Int {
	d := new(big.Int)
	d.Sub(max, min)
	d.Div(d, big.NewInt(2))

	r := new(big.Int).Set(min)
	r.Add(r, big.NewInt(rand.Int63n(d.Int64())))
	r.Add(r, big.NewInt(rand.Int63n(d.Int64())))

	return r
}
