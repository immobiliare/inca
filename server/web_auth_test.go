package server

import (
	"io"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/util"
	"github.com/matryer/is"
)

const (
	testingToken = "awesometoken"
)

func testAuthApp(t *testing.T) *Inca {
	var (
		app     = testApp(t)
		authApp = &Inca{}
	)

	*authApp = *app
	authApp.acl = map[string][]string{
		testingToken: {`.*.domain.tld`},
	}
	return authApp
}

func TestServerWebAuthLoginView(t *testing.T) {
	t.Parallel()

	var (
		app  = testAuthApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/web/login", nil),
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

func TestServerWebAuthLogin(t *testing.T) {
	t.Parallel()

	var (
		app  = testAuthApp(t)
		test = is.New(t)
	)

	form := url.Values{}
	form.Add("token", testingToken)
	request := httptest.NewRequest("POST", "/web/login", strings.NewReader(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	test.NoErr(err)

	body, err := io.ReadAll(response.Body)
	test.NoErr(err)
	defer func() {
		if err := response.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	test.True(util.IsValidHTML(body))
}

func TestServerWebAuthLogout(t *testing.T) {
	t.Parallel()

	var (
		app  = testAuthApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/web/logout", nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusFound)

	body, err := io.ReadAll(response.Body)
	test.NoErr(err)
	defer func() {
		if err := response.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	test.True(util.IsValidHTML(body))
}
