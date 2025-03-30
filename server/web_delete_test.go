package server

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/util"
	"github.com/matryer/is"
)

func TestServerWebDelete(t *testing.T) {
	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("POST", fmt.Sprintf("/web/%s/delete", testingCADomain), nil),
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
	test.True(strings.Contains(string(body), "not found"))

	response, err = app.Test(
		httptest.NewRequest("GET", "/"+testingCADomain+"?algo="+testingCAAlgorithm, nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	response, err = app.Test(
		httptest.NewRequest("POST", fmt.Sprintf("/web/%s/delete", testingCADomain), nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	body, err = io.ReadAll(response.Body)
	test.NoErr(err)
	defer func() {
		if err := response.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	test.True(util.IsValidHTML(body))
	test.True(!strings.Contains(string(body), "not found"))

}
