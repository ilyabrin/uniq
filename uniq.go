package uniq

import (
	"crypto/rand"
	"encoding/binary"
	mathrand "math/rand"
	"sync"
)

const (
	letters          = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers          = "0123456789"
)

// randSource is a process-wide source of randomness. It is seeded once
// from crypto/rand, so two calls to New in the same nanosecond no longer
// produce identical IDs.
var (
	randMu     sync.Mutex
	randSource = mathrand.New(mathrand.NewSource(cryptoSeed()))
)

// cryptoSeed reads a random int64 from crypto/rand. It panics only if the
// OS entropy source is unavailable, which is unrecoverable anyway.
func cryptoSeed() int64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic("uniq: failed to read crypto/rand: " + err.Error())
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}

func generateRandomStringWithUpperCase(seed *mathrand.Rand, length, numUpperCase int) string {
	if length <= 0 {
		return ""
	}
	if numUpperCase > length {
		numUpperCase = length
	}
	if numUpperCase < 0 {
		numUpperCase = 0
	}

	result := make([]byte, length)

	for i := range result {
		result[i] = letters[seed.Intn(len(letters))]
	}

	// Uppercase exactly numUpperCase distinct positions, chosen via a
	// partial Fisher-Yates shuffle instead of retrying random indexes.
	positions := seed.Perm(length)[:numUpperCase]
	for _, p := range positions {
		result[p] = uppercaseLetters[result[p]-'a']
	}

	return string(result)
}

// New returns a random string of the given length containing exactly
// numUpperCase uppercase letters (capped at length). A non-positive
// length returns an empty string.
func New(length, numUpperCase int) string {
	randMu.Lock()
	defer randMu.Unlock()
	return generateRandomStringWithUpperCase(randSource, length, numUpperCase)
}
