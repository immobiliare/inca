package util

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand/v2"
	"regexp"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

// GenerateRandomString returns a random string of the specified length using characters from a predefined set, including letters, digits, and special symbols. If secure random generation fails, it falls back to a less secure random source.
func GenerateRandomString(length int) string {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			randomIndex = big.NewInt(int64(mrand.IntN(len(charset))))
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result)
}

// RegexesMatch returns true if the payload matches all provided regular expressions.
// If any regex fails to match or an error occurs during matching, the result is false.
func RegexesMatch(payload string, regexes ...string) bool {
	match := true
	for _, filter := range regexes {
		match = match && ErrWrap(false)(regexp.MatchString(filter, payload))
	}
	return match
}

func RegexesAnyMatch(payload string, regexes ...string) bool {
	for _, filter := range regexes {
		if ErrWrap(false)(regexp.MatchString(filter, payload)) {
			return true
		}
	}
	return false
}
