package pki

import (
	"crypto/x509"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

var (
	name      = "domain.tld"
	address   = "127.0.0.1"
	altArray  = []string{name, address}
	altString = fmt.Sprintf("%s,%s", name, address)
	req       = NewRequest(map[string]any{
		"cn":  name,
		"alt": altString,
	})
	testCrt *x509.Certificate
	testKey *Key
)

func newPair(t *testing.T) {
	reqCrt, reqKey, err := New(req)
	is.New(t).NoErr(err)
	testCrt = reqCrt
	testKey = reqKey
}

func crt(t *testing.T) *x509.Certificate {
	if testCrt == nil {
		newPair(t)
	}
	return testCrt
}

func key(t *testing.T) *Key {
	if testKey == nil {
		newPair(t)
	}
	return testKey
}

func TestPkiCrtParseBytes(t *testing.T) {
	var (
		crt  = crt(t)
		key  = key(t)
		test = is.New(t)
	)

	wrap, err := WrapCrt(crt, key, crt, key)
	test.NoErr(err)

	_, err = ParseBytes(ExportBytes(wrap))
	test.NoErr(err)
}
func TestPkiCrtParseKeyPairBytes(t *testing.T) {
	var (
		crt  = crt(t)
		key  = key(t)
		test = is.New(t)
	)

	wrapCrt, wrapKey, err := Wrap(crt, key, crt, key)
	test.NoErr(err)

	_, _, err = ParseKeyPairBytes(ExportBytes(wrapCrt), ExportBytes(wrapKey))
	test.NoErr(err)
}

func TestPkiCrtNew(t *testing.T) {
	var (
		crt  = crt(t)
		test = is.New(t)
	)
	test.Equal(crt.Subject.CommonName, req.CN)
	test.Equal(crt.DNSNames[0], name)
	test.Equal(crt.IPAddresses[0].String(), address)
}

func TestPkiCrtIsValidCN(t *testing.T) {
	var (
		test        = is.New(t)
		commonNames = map[string]bool{
			name:                    true,
			"dom-ain.tld":           true,
			"sub.domain.tld":        true,
			"-domain.tld":           false,
			"domain.super-long-tld": false,
		}
	)
	for commonName, valid := range commonNames {
		test.Equal(IsValidCN(commonName), valid)
	}
}

func TestPkiCrtAltNames(t *testing.T) {
	var (
		crt              = crt(t)
		test             = is.New(t)
		names, addresses = AltNames(crt)
	)
	test.Equal(names, []string{name})
	test.Equal(addresses, []string{address})
}

func TestPkiCrtParseAltNames(t *testing.T) {
	var (
		test             = is.New(t)
		names, addresses = ParseAltNames(altArray)
	)
	test.Equal(names, []string{name})
	test.Equal(addresses, []string{address})
}
