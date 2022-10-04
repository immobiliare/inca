package helper

import (
	"testing"

	"github.com/matryer/is"
)

const (
	testingToken = "hello"
)

var (
	testingACL      = map[string][]string{testingToken: nil}
	testingACLEmpty = map[string][]string{}
)

func TestServerHelperIsValidToken(t *testing.T) {
	is.New(t).True(IsValidToken(testingToken, testingACL))
	is.New(t).True(!IsValidToken(testingToken, testingACLEmpty))
}
