package server

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/matryer/is"
	"gitlab.rete.farm/sistemi/inca/util"
)

const (
	testingImportDomain = "foreign.tld"
	testingImportCrt    = `-----BEGIN CERTIFICATE-----
MIIB4zCCAYmgAwIBAgIQOow+W10AzEJL0mVd4zQogTAKBggqhkjOPQQDAjBYMQkw
BwYDVQQGEwAxCTAHBgNVBAgTADEJMAcGA1UEBxMAMQkwBwYDVQQJEwAxCTAHBgNV
BBETADEJMAcGA1UEChMAMRQwEgYDVQQDEwtmb3JlaWduLnRsZDAeFw0yMjEwMDMx
MzE4MTJaFw0yMzExMDQxMzE4MTJaMFgxCTAHBgNVBAYTADEJMAcGA1UECBMAMQkw
BwYDVQQHEwAxCTAHBgNVBAkTADEJMAcGA1UEERMAMQkwBwYDVQQKEwAxFDASBgNV
BAMTC2ZvcmVpZ24udGxkMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEVQcXumlM
nSoZEmU+Yd47agAo76oi9lVrPYaGNNw9PPgOSLa7fnPAXQsq9teEZ5hHnEZnpfbo
jRGkwLU8dVbtb6M1MDMwDgYDVR0PAQH/BAQDAgeAMBMGA1UdJQQMMAoGCCsGAQUF
BwMBMAwGA1UdEwEB/wQCMAAwCgYIKoZIzj0EAwIDSAAwRQIgdcq630BYzpjJwcY8
6K8rMbSu0SjDqRek865/+wQOwKwCIQCqXIKAXlRoFtxX7sMobmZVuNBVz7a/QPpX
rbRKCJDQ9w==
-----END CERTIFICATE-----`
	testingImportKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgC+8lxEP1SSy9P9JO
Ant8Ven+h/a6kICfo/gi++doWBihRANCAARVBxe6aUydKhkSZT5h3jtqACjvqiL2
VWs9hoY03D08+A5Itrt+c8BdCyr214RnmEecRmel9uiNEaTAtTx1Vu1v
-----END PRIVATE KEY-----`
)

func TestServerWebImportView(t *testing.T) {
	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/web/import", nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	body, err := io.ReadAll(response.Body)
	test.NoErr(err)
	defer response.Body.Close()

	test.True(util.IsValidHTML(body))
}

func TestServerWebImport(t *testing.T) {
	var (
		app  = testApp(t)
		test = is.New(t)
	)

	bodyWriter := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyWriter)
	crt, err := writer.CreateFormFile("crt", "crt.pem")
	test.NoErr(err)
	_, err = io.Copy(crt, bytes.NewReader([]byte(testingImportCrt)))
	test.NoErr(err)
	key, err := writer.CreateFormFile("key", "key.pem")
	test.NoErr(err)
	_, err = io.Copy(key, bytes.NewReader([]byte(testingImportKey)))
	test.NoErr(err)
	test.NoErr(writer.Close())

	request := httptest.NewRequest("POST", "/web/import", bytes.NewReader(bodyWriter.Bytes()))
	request.Header.Add("Content-Type", writer.FormDataContentType())

	response, err := app.Test(request)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	body, err := io.ReadAll(response.Body)
	test.NoErr(err)
	defer response.Body.Close()

	test.True(util.IsValidHTML(body))
	test.True(!strings.Contains(string(body), "danger"))

	response, err = app.Test(
		httptest.NewRequest("GET", fmt.Sprintf("/%s", testingImportDomain), nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	body, err = io.ReadAll(response.Body)
	test.NoErr(err)
	defer response.Body.Close()

	test.True(util.IsValidHTML(body))
	test.True(!strings.Contains(string(body), "not found"))

	response, err = app.Test(
		httptest.NewRequest("DELETE", "/"+testingImportDomain, nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)
}
