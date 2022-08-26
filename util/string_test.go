package util

import (
	"testing"

	"github.com/matryer/is"
)

const testingString = "hello"

func TestUtilStringRegexesMatch(t *testing.T) {
	is.New(t).True(RegexesMatch(testingString, ".*"))
}

func TestUtilStringRegexesMultipleMatch(t *testing.T) {
	is.New(t).True(RegexesMatch(testingString, ".*", ".+"))
}

func TestUtilStringRegexesNotMatch(t *testing.T) {
	is.New(t).True(!RegexesMatch(testingString, "^$"))
}

func TestUtilStringRegexesNotMultipleMatch(t *testing.T) {
	is.New(t).True(!RegexesMatch(testingString, ".*", "^$"))
}
