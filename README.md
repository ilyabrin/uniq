# UNIQ

[![Go Tests](https://github.com/ilyabrin/uniq/actions/workflows/test.yml/badge.svg)](https://github.com/ilyabrin/uniq/actions/workflows/test.yml)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/ilyabrin/uniq.svg?style=flat)](https://github.com/ilyabrin/uniq/pulls)
[![PR's Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](http://makeapullrequest.com)
[![GitHub Release](https://img.shields.io/github/release/ilyabrin/uniq.svg?style=flat)](https://github.com/ilyabrin/uniq/releases)

Random ID and token generation for Go: simple IDs, cryptographically
secure tokens, and a configurable generator with functional options.

## Installation

```sh
go get -u github.com/ilyabrin/uniq
```

## Usage

### Secure tokens — `Token`

Cryptographically secure (crypto/rand, no modulo bias), alphabet
`[a-zA-Z0-9]`. Use for password reset tokens, API keys, session IDs,
CSRF tokens, invite links.

```go
token, err := uniq.Token(32)   // "kX9mPq2RtY7wBn4ZcV1sAj8dHf3GuLe0"
key := uniq.MustToken(64)      // panics only if OS entropy fails
```

`Token(32)` carries ~190 bits of entropy — comfortably above the
128 bits recommended by OWASP for session identifiers.

### Configurable generator — `Generate`

Also crypto/rand-backed. Defaults to lowercase letters; combine
functional options:

```go
id, err := uniq.Generate(20)                              // lowercase only
id  = uniq.MustGenerate(20, uniq.WithDigits())            // + digits
id  = uniq.MustGenerate(20, uniq.WithUppercase(5))        // exactly 5 uppercase
id  = uniq.MustGenerate(6,  uniq.WithAlphabet("0123456789")) // SMS code
id  = uniq.MustGenerate(10, uniq.WithAlphabet("abcdefghjkmnpqrstuvwxyz23456789")) // no look-alikes (0O1lI)
id  = uniq.MustGenerate(32, uniq.WithDigits(), uniq.WithUppercase(8)) // combined
```

Good for coupon and invite codes, URL slugs, SMS confirmation codes,
IDs that must match an external system's charset, or strings with a
guaranteed number of uppercase letters.

Invalid input returns an error instead of panicking: negative length,
empty alphabet, alphabet longer than 256 characters.

### Simple IDs — `New`

The original API: lowercase letters with exactly `numUpperCase`
uppercase. Fast, thread-safe, seeded once from crypto/rand — but not
cryptographically secure (math/rand). Use for log request IDs, temp
file suffixes, and other non-security identifiers.

```go
id := uniq.New(10, 3) // e.g. "aKbcDefGhi" — length 10, 3 uppercase
```

Playground: [go.dev/play/p/i6vKfk-evij](https://go.dev/play/p/i6vKfk-evij)

## Choosing a function

| Need | Use |
| --- | --- |
| Security-sensitive token (auth, API key, session) | `Token` |
| Custom charset / digits / exact uppercase count | `Generate` + options |
| Fast non-security ID | `New` |

## Benchmarks

12th Gen Intel i5-12400F, Go 1.25, `go test -bench . -benchmem`:

| Call | ns/op | B/op | allocs/op |
| --- | --- | --- | --- |
| `New(10, 3)` | 152 | 96 | 2 |
| `New(20, 6)` | 264 | 184 | 2 |
| `Token(16)` | 195 | 16 | 1 |
| `Token(32)` | 282 | 32 | 1 |
| `Token(32)` parallel | 63 | 32 | 1 |
| `Generate(10)` | 228 | 48 | 2 |
| `Generate(20, WithDigits())` | 250 | 56 | 2 |
| `Generate(20, WithUppercase(6))` | 2829 | 216 | 3 |

`Token` and `Generate` scale across cores (no shared lock); `New`
serializes on a mutex. `WithUppercase` is the slowest path — it draws
a crypto-random permutation for the uppercase positions.

## License

MIT
