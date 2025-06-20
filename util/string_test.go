package util

import (
	"testing"

	"github.com/matryer/is"
)

const testingString = "hello"

func TestUtilStringRegexesMatch(t *testing.T) {
	t.Parallel()

	is.New(t).True(RegexesMatch(testingString, ".*"))
}

func TestUtilStringRegexesMultipleMatch(t *testing.T) {
	t.Parallel()

	is.New(t).True(RegexesMatch(testingString, ".*", ".+"))
}

func TestUtilStringRegexesNotMatch(t *testing.T) {
	t.Parallel()

	is.New(t).True(!RegexesMatch(testingString, "^$"))
}

func TestUtilStringRegexesNotMultipleMatch(t *testing.T) {
	t.Parallel()

	is.New(t).True(!RegexesMatch(testingString, ".*", "^$"))
}

func TestUtilStringRegexesAnyMatch(t *testing.T) {
	t.Parallel()

	is.New(t).True(RegexesAnyMatch(testingString, ".*", "^$"))
}

func TestUtilStringRegexesNotAnyMatch(t *testing.T) {
	t.Parallel()

	is.New(t).True(!RegexesAnyMatch(testingString, "^$"))
}

func TestGenerateRandomString_Length(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	for _, l := range []int{0, 1, 5, 10, 50, 100} {
		s := GenerateRandomString(l)
		is.Equal(len(s), l)
	}
}

func TestGenerateRandomString_Charset(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	s := GenerateRandomString(1000)
	for _, c := range s {
		found := false
		for _, allowed := range charset {
			if c == allowed {
				found = true
				break
			}
		}
		is.True(found)
	}
}

func TestGenerateRandomString_Uniqueness(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	s1 := GenerateRandomString(32)
	s2 := GenerateRandomString(32)
	is.True(s1 != s2)
}
