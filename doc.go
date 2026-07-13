// Package uniq generates random identifiers and tokens.
//
// It provides three entry points:
//
//   - [Token] and [MustToken]: cryptographically secure strings over
//     [a-zA-Z0-9], suitable for API keys, session and reset tokens.
//   - [Generate] and [MustGenerate]: crypto/rand-backed generator
//     configured with functional options ([WithDigits], [WithUppercase],
//     [WithAlphabet]).
//   - [New]: fast math/rand-based IDs with an exact number of uppercase
//     letters, for non-security uses such as request IDs.
package uniq
