package dataset

import (
	"math/rand"
	"strings"

	"github.com/bxcodec/faker/v3"
)

func Ascii(min, max int) string {
	charset := []rune("abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789" + "~=+%^*/()[]{}/!@#$?| ")

	l := int(rand.Int63n(int64(max-min))) + min

	s := strings.Builder{}
	for i := 0; i < l; i++ {
		idx := rand.Intn(len(charset))
		s.WriteRune(charset[idx])
	}

	return s.String()
}

func WordN(n int) string {
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = faker.Word()
	}
	return strings.Join(words, ", ")
}

func SentenceN(n int) string {
	s := make([]string, n)
	for i := 0; i < n; i++ {
		s[i] = faker.Sentence()
	}
	return strings.Join(s, " ")
}

func ParagraphN(n int) string {
	p := make([]string, n)
	for i := 0; i < n; i++ {
		p[i] = faker.Paragraph()
	}
	return strings.Join(p, " ")
}
