package server

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/pki"
	"github.com/matryer/is"
)

func TestServerCA(t *testing.T) {
	t.Parallel()

	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/ca/local", nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	body, err := io.ReadAll(response.Body)
	test.NoErr(err)
	defer response.Body.Close()

	crt, err := pki.ParseBytes(body)
	test.NoErr(err)
	test.Equal(crt.Subject.CommonName, testingCADomain)
}

func TestServerCANotFound(t *testing.T) {
	t.Parallel()

	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/ca/letsencrypt", nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusNotFound)
}
