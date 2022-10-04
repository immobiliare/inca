package util

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

const (
	testingSliceData = "hello"
)

var (
	testingSlice      = []string{testingSliceData}
	testingSliceEmpty = []string{}
)

func TestUtilSliceEqual(t *testing.T) {
	is.New(t).True(StringSlicesEqual(testingSlice, testingSlice))
}

func TestUtilSliceNotEqual(t *testing.T) {
	is.New(t).True(!StringSlicesEqual(testingSlice, testingSliceEmpty))
}

func TestUtilSliceDistinct(t *testing.T) {
	is.New(t).Equal(StringSliceDistinct(testingSlice), testingSlice)
}

func TestUtilSliceDistinctDouble(t *testing.T) {
	is.New(t).Equal(StringSliceDistinct(append(testingSlice, testingSliceData)), testingSlice)
}

func TestUtilSliceContains(t *testing.T) {
	is.New(t).True(StringSliceContains(testingSlice, testingSliceData))
}

func TestUtilSliceNotContains(t *testing.T) {
	is.New(t).True(!StringSliceContains(testingSlice, fmt.Sprintf("%s.world", testingSliceData)))
}
