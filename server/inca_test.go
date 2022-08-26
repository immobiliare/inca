package server

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

const (
	testingCADomain   = "domain.tld"
	testingConfigPath = "/tmp/.testServerInca.yml"
	testingCACrtPath  = "/tmp/.testServerInca.crt.pem"
	testingCAKeyPath  = "/tmp/.testServerInca.key.pem"
	testingConfig     = `storage:
    type: fs
    path: /tmp
providers:
    - type: local
      crt: ` + testingCACrtPath + `
      key: ` + testingCAKeyPath + `
`
	testingCACrt = `-----BEGIN CERTIFICATE-----
MIICBDCCAamgAwIBAgIQN4zDrydkUQroRxbD3v1WUjAKBggqhkjOPQQDAjBXMQkw
BwYDVQQGEwAxCTAHBgNVBAgTADEJMAcGA1UEBxMAMQkwBwYDVQQJEwAxCTAHBgNV
BBETADEJMAcGA1UEChMAMRMwEQYDVQQDEwpkb21haW4udGxkMB4XDTIyMDgyNjEz
NDIwMFoXDTIzMDkyNzEzNDIwMFowVzEJMAcGA1UEBhMAMQkwBwYDVQQIEwAxCTAH
BgNVBAcTADEJMAcGA1UECRMAMQkwBwYDVQQREwAxCTAHBgNVBAoTADETMBEGA1UE
AxMKZG9tYWluLnRsZDBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABGXWem02EVha
1w9UcPz13NB1uVG0jsFr1KdYO1RDujBQZ6ga0dEbuSxkRe6qlh4QXUhtXqZkV73a
SXXJBbOLLWejVzBVMA4GA1UdDwEB/wQEAwIChDATBgNVHSUEDDAKBggrBgEFBQcD
ATAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBRtpgtnscJakHV46qHAHS/PfA9S
8TAKBggqhkjOPQQDAgNJADBGAiEAwrgfXf9UoLMbpZ7HnjCsu4/33vMhGLFMuEeX
VGILY+0CIQDssfxoYJ8khXDL4y72DsNlcM2JZ8n8xCy1duGeXwGjyg==
-----END CERTIFICATE-----`
	testingCAKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg2Z2ABjYeQV0twHeA
zPFOeNIjv75tBZzY2oROCALQWRGhRANCAARl1nptNhFYWtcPVHD89dzQdblRtI7B
a9SnWDtUQ7owUGeoGtHRG7ksZEXuqpYeEF1IbV6mZFe92kl1yQWziy1n
-----END PRIVATE KEY-----`
)

var testingApp *Inca

func testApp(t *testing.T) *Inca {
	if testingApp == nil {
		test := is.New(t)
		test.NoErr(
			os.WriteFile(testingConfigPath, []byte(testingConfig), 0644))
		defer os.Remove(testingConfigPath)
		test.NoErr(
			os.WriteFile(testingCACrtPath, []byte(testingCACrt), 0644))
		defer os.Remove(testingCACrtPath)
		test.NoErr(
			os.WriteFile(testingCAKeyPath, []byte(testingCAKey), 0644))
		defer os.Remove(testingCAKeyPath)

		testApp, err := Spinup(testingConfigPath)
		is.New(t).NoErr(err)
		testingApp = testApp
	}
	return testingApp
}
