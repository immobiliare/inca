package server

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/matryer/is"
)

const (
	testingCADomain    = "domain.tld"
	testingCAAlgorithm = "ecdsa"
	testingConfigPath  = ".testServerInca.yml"
	testingCACrtPath   = ".testServerInca.crt.pem"
	testingCAKeyPath   = ".testServerInca.key.pem"
	testingConfig      = `bind: :65535
storage:
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
		defer func() {
			if err := os.Remove(testingConfigPath); err != nil {
				t.Logf("Failed to remove test config file: %v", err)
			}
		}()

		test.NoErr(
			os.WriteFile(testingCACrtPath, []byte(testingCACrt), 0644))
		defer func() {
			if err := os.Remove(testingCACrtPath); err != nil {
				t.Logf("Failed to remove test CA certificate file: %v", err)
			}
		}()

		test.NoErr(
			os.WriteFile(testingCAKeyPath, []byte(testingCAKey), 0644))
		defer func() {
			if err := os.Remove(testingCAKeyPath); err != nil {
				t.Logf("Failed to remove test CA key file: %v", err)
			}
		}()

		testApp, err := Spinup(testingConfigPath)
		is.New(t).NoErr(err)
		testingApp = testApp
	}
	return testingApp
}

func TestStatic(t *testing.T) {
	t.Parallel()

	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/static/favicon.ico", nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	response, err = app.Test(
		httptest.NewRequest("GET", "/static/style.css", nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)
}
