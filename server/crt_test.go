package server

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/pki"
	"github.com/matryer/is"
)

func TestServerCrt(t *testing.T) {
	t.Parallel()

	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/"+testingCADomain+"?algo="+testingCAAlgorithm, nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	body, err := io.ReadAll(response.Body)
	test.NoErr(err)
	defer func() {
		if err := response.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	crt, err := pki.ParseBytes(body)
	test.NoErr(err)
	test.Equal(crt.Subject.CommonName, testingCADomain)
	test.Equal(crt.Issuer.CommonName, testingCADomain)

	response, err = app.Test(
		httptest.NewRequest("DELETE", "/"+testingCADomain, nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)
}

func TestServerCrtBad(t *testing.T) {
	t.Parallel()

	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/domain2.tld", nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusBadRequest)
}
