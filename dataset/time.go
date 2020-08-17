package dataset

import (
	"math/rand"
	"time"
)

func DateTime() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2038, 1, 19, 3, 14, 17, 0, time.UTC).Unix()
	d := max - min

	return time.Unix(rand.Int63n(d)+min, 0)
}
