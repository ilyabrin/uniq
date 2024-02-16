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

// func generateRandomStringWithUpperCase(seed *rand.Rand, length, numUpperCase int) string {

// 	result := make([]byte, length)
// 	countUpperCase := 0 // counter for uppercase chars

// 	// generate random string with lower case
// 	for i := range result {
// 		result[i] = letters[seed.Intn(len(letters))]
// 	}

// 	// change case only for lowercase letters
// 	for i := 0; i < numUpperCase; i++ {
// 		randIndex := seed.Intn(length)
// 		if result[randIndex] >= 'a' && result[randIndex] <= 'z' {
// 			result[randIndex] = uppercaseLetters[result[randIndex]-'a']
// 		}
// 	}

// 	// count uppercase letters
// 	for i := 0; i < len(result); i++ {
// 		if result[i] >= 'A' && result[i] <= 'Z' {
// 			countUpperCase++
// 		}
// 	}

// 	// change lower case to upper case in place
// 	for countUpperCase < numUpperCase {
// 		randIndex := seed.Intn(length)
// 		if result[randIndex] >= 'a' && result[randIndex] <= 'z' {
// 			result[randIndex] = uppercaseLetters[result[randIndex]-'a']
// 			countUpperCase++
// 		}
// 	}

// 	return string(result)
// }

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
