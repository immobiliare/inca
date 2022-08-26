package pki

import (
	"testing"

	"github.com/matryer/is"
)

func TestPkiBundleWrapCrt(t *testing.T) {
	var (
		crt  = testCrt(t)
		key  = key(t)
		test = is.New(t)
	)

	wrap, err := WrapCrt(crt, key, crt, key)
	test.NoErr(err)
	test.True(len(ExportBytes(wrap)) > 0)
}

func TestPkiBundleWrapKey(t *testing.T) {
	var (
		key  = key(t)
		test = is.New(t)
	)

	wrap, err := WrapKey(key)
	test.NoErr(err)
	test.True(len(ExportBytes(wrap)) > 0)
}

func TestPkiBundleWrap(t *testing.T) {
	var (
		crt  = testCrt(t)
		key  = key(t)
		test = is.New(t)
	)

	wrapCrt, wrapKey, err := Wrap(crt, key, crt, key)
	test.NoErr(err)
	test.True(len(ExportBytes(wrapCrt)) > 0)
	test.True(len(ExportBytes(wrapKey)) > 0)
}
