package uniq

import (
	"math/rand"
	"time"
)

const (
	letters          = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers          = "0123456789"
)

func generateRandomStringWithUpperCase(seed *rand.Rand, length, numUpperCase int) string {

	if numUpperCase > length {
		numUpperCase = length
	}

	result := make([]byte, length)

	for i := range result {
		result[i] = letters[seed.Intn(len(letters))]
	}

	for i := 0; i < numUpperCase; {
		randIndex := seed.Intn(length)
		if result[randIndex] >= 'a' && result[randIndex] <= 'z' {
			result[randIndex] = uppercaseLetters[result[randIndex]-'a']
			i++
		}
	}

	return string(result)
}

func New(length, size int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return generateRandomStringWithUpperCase(seededRand, length, size)
}
