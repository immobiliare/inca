package helper

import (
	"testing"

	"github.com/matryer/is"
)

const (
	testingToken = "hello"
)

var (
	testingACLValid   = map[string][]string{testingToken: nil}
	testingACLInvalid = map[string][]string{"world": nil}
	testingACLEmpty   = map[string][]string{}
)

func TestServerHelperIsValidToken(t *testing.T) {
	t.Parallel()

	is.New(t).True(IsValidToken(testingToken, testingACLValid))
	is.New(t).True(!IsValidToken(testingToken, testingACLInvalid))
	is.New(t).True(!IsValidToken(testingToken, testingACLEmpty))
}
