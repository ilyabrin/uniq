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
		positions, err := randomPerm(length)
		if err != nil {
			return "", err
		}
		lower := make([]byte, numUpper)
		if err := fillRandom(lower, letters); err != nil {
			return "", err
		}
		for i, p := range positions[:numUpper] {
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

// randomPerm returns a random permutation of [0, n) built with a
// crypto/rand-backed Fisher-Yates shuffle.
func randomPerm(n int) ([]int, error) {
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}
	idx := make([]byte, 1)
	for i := n - 1; i > 0; i-- {
		// Uniform j in [0, i] via rejection sampling over one byte;
		// fine for n <= 256, and split reads keep it general enough
		// for typical ID lengths.
		bound := i + 1
		if bound > 256 {
			return nil, errors.New("uniq: length must not exceed 256 when using WithUppercase")
		}
		limit := 256 / bound * bound
		for {
			if _, err := rand.Read(idx); err != nil {
				return nil, err
			}
			if int(idx[0]) < limit {
				j := int(idx[0]) % bound
				perm[i], perm[j] = perm[j], perm[i]
				break
			}
		}
	}
	return perm, nil
}
