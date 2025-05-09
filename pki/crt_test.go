package pki

import (
	"crypto/x509"
	"fmt"
	"os"
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
	reqCrt, reqKey, err := New(testingReq, false)
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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

	var (
		crt  = testCrt(t)
		test = is.New(t)
	)
	test.Equal(crt.Subject.CommonName, testingReq.CN)
	test.Equal(crt.DNSNames[0], testingName)
	test.Equal(crt.IPAddresses[0].String(), testingAddress)
}

func TestPkiCrtIsValidCN(t *testing.T) {
	t.Parallel()

	var (
		test        = is.New(t)
		commonNames = map[string]bool{
			testingName:             true,
			"dom-ain.tld":           true,
			"sub.domain.tld":        true,
			"*.domain.tld":          true,
			"-domain.tld":           false,
			"domain.super-long-tld": false,
		}
	)
	for commonName, valid := range commonNames {
		test.Equal(IsValidCN(commonName), valid)
	}
}

func TestPkiCrtAltNames(t *testing.T) {
	t.Parallel()

	var (
		crt              = testCrt(t)
		test             = is.New(t)
		names, addresses = AltNames(crt)
	)
	test.Equal(names, []string{testingName})
	test.Equal(addresses, []string{testingAddress})
}

func TestPkiCrtParseAltNames(t *testing.T) {
	t.Parallel()

	var (
		test             = is.New(t)
		names, addresses = ParseAltNames(testingAltArray)
	)
	test.Equal(names, []string{testingName})
	test.Equal(addresses, []string{testingAddress})
}

func TestPkiCrtParse(t *testing.T) {
	t.Parallel()

	var (
		crt  = testCrt(t)
		key  = key(t)
		test = is.New(t)
	)

	wrap, err := WrapCrt(crt, key, crt, key)
	test.NoErr(err)

	f, err := os.CreateTemp("", "test-cert-*.pem")
	test.NoErr(err)
	defer func() {
		if err := os.Remove(f.Name()); err != nil {
			t.Logf("Failed to remove temporary file: %v", err)
		}
	}()

	err = os.WriteFile(f.Name(), ExportBytes(wrap), 0644)
	test.NoErr(err)

	parsedCrt, err := Parse(f.Name())
	test.NoErr(err)
	test.Equal(parsedCrt.Subject.CommonName, crt.Subject.CommonName)

	_, err = Parse("non-existent-file.pem")
	test.True(err != nil)
}

func TestPkiCrtParseKeyPair(t *testing.T) {
	t.Parallel()

	var (
		crt  = testCrt(t)
		key  = key(t)
		test = is.New(t)
	)

	wrapCrt, wrapKey, err := Wrap(crt, key, crt, key)
	test.NoErr(err)

	crtFile, err := os.CreateTemp("", "test-cert-*.pem")
	test.NoErr(err)
	defer func() {
		if err := os.Remove(crtFile.Name()); err != nil {
			t.Logf("Failed to remove certificate file: %v", err)
		}
	}()

	keyFile, err := os.CreateTemp("", "test-key-*.pem")
	test.NoErr(err)
	defer func() {
		if err := os.Remove(keyFile.Name()); err != nil {
			t.Logf("Failed to remove key file: %v", err)
		}
	}()

	err = os.WriteFile(crtFile.Name(), ExportBytes(wrapCrt), 0644)
	test.NoErr(err)

	err = os.WriteFile(keyFile.Name(), ExportBytes(wrapKey), 0644)
	test.NoErr(err)

	parsedCrt, parsedKey, err := ParseKeyPair(crtFile.Name(), keyFile.Name())
	test.NoErr(err)
	test.Equal(parsedCrt.Subject.CommonName, crt.Subject.CommonName)
	test.True(parsedKey != nil)

	_, _, err = ParseKeyPair("non-existent-cert.pem", "non-existent-key.pem")
	test.True(err != nil)
}
