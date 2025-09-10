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
	// Don't run in parallel to avoid state conflicts
	var (
		app  = testApp(t)
		test = is.New(t)
	)

	// Clean up any existing certificates first
	if resp, err := app.Test(httptest.NewRequest("DELETE", "/"+testingCADomain, nil)); err != nil {
		t.Logf("pre-cleanup DELETE failed: %v", err)
	} else if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		if cerr := resp.Body.Close(); cerr != nil {
			t.Logf("Failed to close response body: %v", cerr)
		}
	}

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
	// Don't run in parallel to avoid state conflicts
	var (
		app  = testApp(t)
		test = is.New(t)
	)

	// Clean up any existing certificates first
	if resp, err := app.Test(httptest.NewRequest("DELETE", "/domain2.tld", nil)); err != nil {
		t.Logf("pre-cleanup DELETE failed: %v", err)
	} else if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		if cerr := resp.Body.Close(); cerr != nil {
			t.Logf("Failed to close response body: %v", cerr)
		}
	}

	response, err := app.Test(
		httptest.NewRequest("GET", "/domain2.tld", nil),
	)
	test.NoErr(err)
	// This should return 400 (Bad Request) for domain2.tld since no provider is found
	test.Equal(response.StatusCode, fiber.StatusBadRequest)
}
