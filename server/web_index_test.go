package server

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/util"
	"github.com/matryer/is"
)

func TestServerWebIndex(t *testing.T) {
	t.Parallel()

	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/", nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusMovedPermanently)

	response, err = app.Test(
		httptest.NewRequest("GET", "/web", nil),
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

	test.True(util.IsValidHTML(body))
}
