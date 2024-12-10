package cmd

import (
	"encoding/json"
	"testing"

	"github.com/immobiliare/inca/pki"
	"github.com/matryer/is"
)

const (
	testingName          = "domain.tld"
	testingOrganization  = "organization"
	testingCountry       = "country"
	testingProvince      = "province"
	testingLocality      = "locality"
	testingStreetAddress = "street address"
	testingPostalCode    = "postal code"
)

func TestCmdGen(t *testing.T) {
	t.Parallel()

	test := is.New(t)
	output, err := testExecute(cmdGen,
		"gen",
		"--ca",
		"--name", testingName,
		"--output", "-",
		"--encode", "json",
		"--organization", testingOrganization,
		"--country", testingCountry,
		"--province", testingProvince,
		"--locality", testingLocality,
		"--street-address", testingStreetAddress,
		"--postal-code", testingPostalCode,
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
	test.Equal(crt.Subject.CommonName, testingName)
	test.Equal(crt.Subject.Organization[0], testingOrganization)
	test.Equal(crt.Subject.Country[0], testingCountry)
	test.Equal(crt.Subject.Province[0], testingProvince)
	test.Equal(crt.Subject.Locality[0], testingLocality)
	test.Equal(crt.Subject.StreetAddress[0], testingStreetAddress)
	test.Equal(crt.Subject.PostalCode[0], testingPostalCode)
}
