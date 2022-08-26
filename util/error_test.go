package util

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

func TestUtilErrorWrapTrue(t *testing.T) {
	is.New(t).True(ErrWrap(true)(func() (bool, error) { return false, errors.New("test") }()))
}

func TestUtilErrorWrapFalse(t *testing.T) {
	is.New(t).True(!ErrWrap(true)(func() (bool, error) { return false, nil }()))
}
