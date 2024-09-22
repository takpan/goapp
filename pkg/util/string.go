package util

import (
	"math/rand"
)

const hexCharBytes = "0123456789ABCDEF" // Possible charaters that a Hex string may contain

var randx = rand.NewSource(42)

// RandString returns a random hex string of length n.
func RandString(n int) string {
	const (
		hexCharIdxBits = 4                     // 4 bits to represent a Hex character index (16 values total)
		hexCharIdxMask = 1<<hexCharIdxBits - 1 // All 1-bits, as many as hexCharIdxBits
		hexCharIdxMax  = 63 / hexCharIdxBits   // # of hex char indices fitting in 63 bits
	)

	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for hexCharIdxMax characters!
	for i, cache, remain := n-1, randx.Int63(), hexCharIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randx.Int63(), hexCharIdxMax
		}
		if idx := int(cache & hexCharIdxMask); idx < len(hexCharBytes) {
			b[i] = hexCharBytes[idx]
			i--
		}
		cache >>= hexCharIdxBits
		remain--
	}

	return string(b)
}
