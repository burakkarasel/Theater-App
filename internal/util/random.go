package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// init runs as program starts and enables true randomness
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random int between given arguments
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string with given char count
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomName returns a 6 random digits string
func RandomName() string {
	return RandomString(6)
}
