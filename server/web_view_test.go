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

func TestServerWebView(t *testing.T) {
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

	response, err = app.Test(
		httptest.NewRequest("GET", fmt.Sprintf("/web/%s", testingCADomain), nil),
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
	test.True(!strings.Contains(string(body), "not found"))
}
