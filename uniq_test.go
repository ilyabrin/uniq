package uniq

import (
	"math/rand"
	"testing"
	"time"
)

func TestGenerateRandomStringWithUpperCase(t *testing.T) {

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []struct {
		length       int
		numUpperCase int
	}{
		{10, 3},
		{15, 5},
		{20, 7},
		{30, 12},
		{50, 20},
	}

	for _, tt := range tests {
		result := generateRandomStringWithUpperCase(seed, tt.length, tt.numUpperCase)

		t.Log(result, tt.length, len(result))

		if len(result) != tt.length {
			t.Errorf("Generated string %s length is incorrect. Expected: %d, Got: %d", result, tt.length, len(result))
		}

		numUpperCase := 0
		for _, char := range result {
			if char >= 'A' && char <= 'Z' {
				numUpperCase++
			}
		}

		if numUpperCase != tt.numUpperCase {
			t.Errorf("Number of uppercase in %s characters is incorrect. Expected: %d, Got: %d", result, tt.numUpperCase, numUpperCase)
		}
	}
}
