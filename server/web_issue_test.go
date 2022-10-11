package server

import (
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/util"
	"github.com/matryer/is"
)

func TestServerWebIssueView(t *testing.T) {
	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/web/issue", nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	body, err := io.ReadAll(response.Body)
	test.NoErr(err)
	defer response.Body.Close()

	test.True(util.IsValidHTML(body))
}

func TestServerWebIssue(t *testing.T) {
	var (
		app        = testApp(t)
		test       = is.New(t)
		testDomain = fmt.Sprintf("testserverwebissue.%s", testingCADomain)
	)

	form := url.Values{}
	form.Add("alt", testDomain)
	form.Add("algo", testingCAAlgorithm)
	request := httptest.NewRequest("POST", "/web/issue", strings.NewReader(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	body, err := io.ReadAll(response.Body)
	test.NoErr(err)
	defer response.Body.Close()

	test.True(util.IsValidHTML(body))
	test.True(strings.Contains(string(body), testDomain))

	response, err = app.Test(
		httptest.NewRequest("DELETE", "/"+testDomain, nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)
}
