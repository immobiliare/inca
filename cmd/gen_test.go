package cmd

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"gitlab.rete.farm/sistemi/inca/pki"
)

const (
	name          = "domain.tld"
	organization  = "organization"
	country       = "country"
	province      = "province"
	locality      = "locality"
	streetAddress = "street address"
	postalCode    = "postal code"
)

func TestCmdGen(t *testing.T) {
	test := is.New(t)
	output, err := mockExecute(cmdGen,
		"gen",
		"--ca",
		"--name", name,
		"--output", "-",
		"--encode", "json",
		"--organization", organization,
		"--country", country,
		"--province", province,
		"--locality", locality,
		"--street-address", streetAddress,
		"--postal-code", postalCode,
	)
	test.NoErr(err)

	outputJson := make(map[string]interface{})
	err = json.Unmarshal([]byte(output), &outputJson)
	test.NoErr(err)

	crtBundle, crtOk := outputJson["crt"]
	keyBundle, keyOk := outputJson["key"]
	test.True(crtOk)
	test.True(keyOk)

	crt, _, err := pki.ParseKeyPairBytes(
		[]byte(crtBundle.(string)),
		[]byte(keyBundle.(string)),
	)
	test.NoErr(err)
	test.True(crt.IsCA)
	test.Equal(crt.Subject.CommonName, name)
	test.Equal(crt.Subject.Organization[0], organization)
	test.Equal(crt.Subject.Country[0], country)
	test.Equal(crt.Subject.Province[0], province)
	test.Equal(crt.Subject.Locality[0], locality)
	test.Equal(crt.Subject.StreetAddress[0], streetAddress)
	test.Equal(crt.Subject.PostalCode[0], postalCode)
}
