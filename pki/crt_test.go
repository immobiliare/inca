package pki

import (
	"crypto/x509"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

var (
	testingName      = "domain.tld"
	testingAddress   = "127.0.0.1"
	testingAltArray  = []string{testingName, testingAddress}
	testingAltString = fmt.Sprintf("%s,%s", testingName, testingAddress)
	testingReq       = NewRequest(map[string]any{
		"cn":  testingName,
		"alt": testingAltString,
	})
	testingCrt *x509.Certificate
	testingKey *Key
)

func testNewPair(t *testing.T) {
	reqCrt, reqKey, err := New(testingReq)
	is.New(t).NoErr(err)
	testingCrt = reqCrt
	testingKey = reqKey
}

func testCrt(t *testing.T) *x509.Certificate {
	if testingCrt == nil {
		testNewPair(t)
	}
	return testingCrt
}

func key(t *testing.T) *Key {
	if testingKey == nil {
		testNewPair(t)
	}
	return testingKey
}

func TestPkiCrtParseBytes(t *testing.T) {
	var (
		crt  = testCrt(t)
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
		crt  = testCrt(t)
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
		crt  = testCrt(t)
		test = is.New(t)
	)
	test.Equal(crt.Subject.CommonName, testingReq.CN)
	test.Equal(crt.DNSNames[0], testingName)
	test.Equal(crt.IPAddresses[0].String(), testingAddress)
}

func TestPkiCrtIsValidCN(t *testing.T) {
	var (
		test        = is.New(t)
		commonNames = map[string]bool{
			testingName:             true,
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
		crt              = testCrt(t)
		test             = is.New(t)
		names, addresses = AltNames(crt)
	)
	test.Equal(names, []string{testingName})
	test.Equal(addresses, []string{testingAddress})
}

func TestPkiCrtParseAltNames(t *testing.T) {
	var (
		test             = is.New(t)
		names, addresses = ParseAltNames(testingAltArray)
	)
	test.Equal(names, []string{testingName})
	test.Equal(addresses, []string{testingAddress})
}
