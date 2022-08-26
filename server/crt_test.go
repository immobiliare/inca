package server

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/matryer/is"
	"gitlab.rete.farm/sistemi/inca/pki"
)

func TestServerCrt(t *testing.T) {
	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/"+testingCADomain, nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	body, err := io.ReadAll(response.Body)
	test.NoErr(err)
	defer response.Body.Close()

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
