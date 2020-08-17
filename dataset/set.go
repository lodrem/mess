package dataset

import (
	"math/rand"
	"strings"
)

func Set(options []string) string {
	var selected []string

	idx := 0
	for {
		idx = idx + rand.Intn(len(options))
		if idx >= len(options) {
			break
		}
		selected = append(selected, options[idx])
	}
	return strings.Join(selected, ",")
}

func Enum(options []string) string {
	return options[rand.Intn(len(options))]
}
