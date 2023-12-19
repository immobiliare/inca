package server

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/util"
	"github.com/matryer/is"
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
		app                 = testApp(t)
		test                = is.New(t)
		testingImportDomain = "foreign.tld"
	)

	testingImportCrt, testingImportKey, err := generateKeyPair(testingImportDomain)
	if err != nil {
		t.Error(err.Error())
	}

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

func generateKeyPair(domain string) (publicX509 string, privateX509 string, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return
	}

	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return
	}

	privateX509 = string(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}))

	publicKey := &privateKey.PublicKey

	subject := pkix.Name{
		CommonName: domain,
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(30) * 24 * time.Hour)

	certTemplate := x509.Certificate{
		Subject:      subject,
		SerialNumber: big.NewInt(1),
		NotBefore:    notBefore,
		NotAfter:     notAfter,
	}

	certDER, err := x509.CreateCertificate(
		rand.Reader,
		&certTemplate,
		&certTemplate,
		publicKey,
		privateKey,
	)
	if err != nil {
		return
	}

	publicX509 = string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}))
	return
}
