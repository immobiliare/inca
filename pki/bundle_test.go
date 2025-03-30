package pki

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestPkiBundleWrapCrt(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	var (
		key  = key(t)
		test = is.New(t)
	)

	wrap, err := WrapKey(key)
	test.NoErr(err)
	test.True(len(ExportBytes(wrap)) > 0)
}

func TestPkiBundleWrap(t *testing.T) {
	t.Parallel()

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

func TestPkiBundleExport(t *testing.T) {
	t.Parallel()

	var (
		crt  = testCrt(t)
		key  = key(t)
		test = is.New(t)
	)

	crtBlock, err := WrapCrt(crt, key, crt, key)
	test.NoErr(err)
	test.NoErr(Export(crtBlock, "test.crt"))
	defer func() {
		if err := os.Remove("test.crt"); err != nil {
			t.Logf("Failed to remove test file: %v", err)
		}
	}()

	keyBlock, err := WrapKey(key)
	test.NoErr(err)
	test.NoErr(Export(keyBlock, "test.key"))
	defer func() {
		if err := os.Remove("test.key"); err != nil {
			t.Logf("Failed to remove test file: %v", err)
		}
	}()

	crtInfo, err := os.Stat("test.crt")
	test.NoErr(err)
	test.Equal(crtInfo.Mode(), os.FileMode(0644))

	keyInfo, err := os.Stat("test.key")
	test.NoErr(err)
	test.Equal(keyInfo.Mode(), os.FileMode(0600))
}
