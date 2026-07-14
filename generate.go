package uniq

import (
	"crypto/rand"
	"errors"
	"fmt"
)

// config holds settings assembled from Options.
type config struct {
	alphabet     string
	numUpperCase int // -1 means "not constrained"
	digits       bool
}

// Option configures Generate.
type Option func(*config)

// WithUppercase makes the result contain exactly n uppercase letters.
// The remaining positions are drawn from the configured alphabet.
func WithUppercase(n int) Option {
	return func(c *config) { c.numUpperCase = n }
}

// WithDigits adds 0-9 to the default alphabet.
func WithDigits() Option {
	return func(c *config) { c.digits = true }
}

// WithAlphabet replaces the alphabet entirely with the given characters.
// It overrides WithDigits.
func WithAlphabet(alphabet string) Option {
	return func(c *config) { c.alphabet = alphabet }
}

// Generate returns a cryptographically secure random string of the given
// length. By default it draws from lowercase letters; combine Options to
// add digits, force an exact number of uppercase letters, or supply a
// custom alphabet.
func Generate(length int, opts ...Option) (string, error) {
	if length < 0 {
		return "", fmt.Errorf("uniq: length must be >= 0, got %d", length)
	}
	if length == 0 {
		return "", nil
	}

	cfg := config{alphabet: letters, numUpperCase: -1}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.digits && cfg.alphabet == letters {
		cfg.alphabet = letters + numbers
	}
	if len(cfg.alphabet) == 0 {
		return "", errors.New("uniq: alphabet must not be empty")
	}
	numUpper := cfg.numUpperCase
	if numUpper > length {
		numUpper = length
	}

	result := make([]byte, length)
	if err := fillRandom(result, cfg.alphabet); err != nil {
		return "", err
	}

	if numUpper >= 0 {
		// Lowercase positions come from the alphabet as usual; then
		// exactly numUpper distinct positions are drawn from a-z and
		// uppercased, replacing whatever was there.
		positions, err := randomPositions(length, numUpper)
		if err != nil {
			return "", err
		}
		lower := make([]byte, numUpper)
		if err := fillRandom(lower, letters); err != nil {
			return "", err
		}
		for i, p := range positions {
			result[p] = uppercaseLetters[lower[i]-'a']
		}
	}

	return string(result), nil
}

// MustGenerate is like Generate but panics on error. Use it when the
// arguments are compile-time constants and entropy failure is fatal.
func MustGenerate(length int, opts ...Option) string {
	s, err := Generate(length, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// fillRandom fills dst with uniformly distributed characters from
// alphabet using crypto/rand and rejection sampling (no modulo bias).
func fillRandom(dst []byte, alphabet string) error {
	n := len(alphabet)
	if n > 256 {
		return errors.New("uniq: alphabet must not exceed 256 characters")
	}
	// Largest multiple of n that fits in a byte; values above it are
	// re-drawn to keep the distribution uniform.
	limit := byte(256 / n * n)

	buf := make([]byte, len(dst))
	filled := 0
	for filled < len(dst) {
		if _, err := rand.Read(buf); err != nil {
			return err
		}
		for _, b := range buf {
			if limit != 0 && b >= limit {
				continue
			}
			dst[filled] = alphabet[int(b)%n]
			filled++
			if filled == len(dst) {
				break
			}
		}
	}
	return nil
}

// randomPositions returns k distinct uniformly random indexes in [0, n)
// via a partial crypto/rand-backed Fisher-Yates shuffle: only the first
// k swaps of the full shuffle are performed, and entropy is read in
// batches instead of one byte per swap.
func randomPositions(n, k int) ([]int, error) {
	if n > 256 {
		return nil, errors.New("uniq: length must not exceed 256 when using WithUppercase")
	}
	if k > n {
		k = n
	}
	if k <= 0 {
		return nil, nil
	}

	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}

	// Batched entropy: refill buf as it drains; rejected bytes just
	// advance the cursor, so on average ~1.2 bytes are used per swap.
	buf := make([]byte, 64)
	pos := len(buf)

	for i := 0; i < k; i++ {
		// Uniform j in [i, n) via rejection sampling over one byte
		// (valid because n <= 256).
		bound := n - i
		limit := 256 / bound * bound
		for {
			if pos == len(buf) {
				if _, err := rand.Read(buf); err != nil {
					return nil, err
				}
				pos = 0
			}
			b := int(buf[pos])
			pos++
			if b < limit {
				j := i + b%bound
				perm[i], perm[j] = perm[j], perm[i]
				break
			}
		}
	}
	return perm[:k], nil
}
