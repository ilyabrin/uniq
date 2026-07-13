package uniq

// tokenAlphabet is the character set for Token: 62 characters, giving
// ~5.95 bits of entropy per character.
const tokenAlphabet = letters + uppercaseLetters + numbers

// Token returns a cryptographically secure random string of the given
// length, built from lowercase letters, uppercase letters and digits.
// Unlike New, it is suitable for security-sensitive values: password
// reset tokens, API keys, session tokens.
//
// It returns an error only if the OS entropy source fails.
func Token(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	result := make([]byte, length)
	if err := fillRandom(result, tokenAlphabet); err != nil {
		return "", err
	}
	return string(result), nil
}

// MustToken is like Token but panics if the OS entropy source fails.
// Convenient for initialization code where such a failure is fatal.
func MustToken(length int) string {
	s, err := Token(length)
	if err != nil {
		panic("uniq: crypto/rand failed: " + err.Error())
	}
	return s
}
