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
