package server

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/matryer/is"
	"gitlab.rete.farm/sistemi/inca/util"
)

func TestServerWebView(t *testing.T) {
	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/"+testingCADomain, nil),
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
	defer response.Body.Close()

	test.True(util.IsValidHTML(body))
	test.True(!strings.Contains(string(body), "not found"))
}
