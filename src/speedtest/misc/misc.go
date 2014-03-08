package misc

import (
	"math/rand"
	"strconv"
)

// shortcut to parse float
func ToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func Urandom(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rand.Int31())
	}

	return b
}
