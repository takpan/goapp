package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	name   string
	strLen int
}

// isStringValid checks whether all characters in testStr are part of the set of validChars
func isStrValid(testStr, validChars string) bool {
	for _, char := range testStr {
		if !strings.Contains(validChars, string(char)) {
			return false
		}
	}
	return true
}

func TestRandString(t *testing.T) {
	var randStringTests = []testData{
		{
			"Empty string", 0,
		},
		{
			"String with 10 characters", 10,
		},
		{
			"String with 1000 characters", 1000,
		},
	}

	for _, tc := range randStringTests {
		t.Run(tc.name, func(t *testing.T) {
			s := RandString(tc.strLen)
			assert.Equal(t, tc.strLen, len(s), "Expected string length to be %d, but got %d", tc.strLen, len(s))
			assert.True(t, isStrValid(s, hexCharBytes), "Generated string contains invalid characters")
		})
	}
}

func BenchmarkRandString(b *testing.B) {
	var randStringBmTests = []testData{
		{
			"String length 10", 10,
		},
		{
			"String length 100", 100,
		},
		{
			"String length 1000", 1000,
		},
	}

	for _, bm := range randStringBmTests {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				RandString(bm.strLen)
			}
		})
	}
}
