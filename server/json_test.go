package server

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	"github.com/matryer/is"
)

func TestServerJSONCrt(t *testing.T) {
	var (
		test = is.New(t)
		crt  = EncodeCrt(&x509.Certificate{
			Subject: pkix.Name{
				CommonName:    testingCADomain,
				Organization:  []string{"organization"},
				Country:       []string{"country"},
				Province:      []string{"province"},
				Locality:      []string{"locality"},
				StreetAddress: []string{"street address"},
				PostalCode:    []string{"postal code"},
			},
		}, nil)
	)
	test.Equal(crt.CN, testingCADomain)
	test.Equal(crt.Organization, "organization")
	test.Equal(crt.Country, "country")
	test.Equal(crt.Province, "province")
	test.Equal(crt.Locality, "locality")
	test.Equal(crt.StreetAddress, "street address")
	test.Equal(crt.PostalCode, "postal code")
	test.Equal(crt.Provider.ID, foreignProviderID)
}
