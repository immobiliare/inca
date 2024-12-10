package pki

import (
	"testing"

	"github.com/matryer/is"
)

func TestPkiKeyNew(t *testing.T) {
	t.Parallel()

	var (
		test = is.New(t)
	)
	for _, algo := range []string{EDDSA, ECDSA, RSA} {
		key, err := newKey(algo)
		test.NoErr(err)
		test.True(key.Public() != nil)
	}
}
