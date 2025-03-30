package server

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/matryer/is"
)

func TestServerEnum(t *testing.T) {
	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/"+testingCADomain+"?algo="+testingCAAlgorithm, nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	response, err = app.Test(
		httptest.NewRequest("GET", "/enum", nil),
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

	bodyJson := make(map[string]interface{})
	err = json.Unmarshal(body, &bodyJson)
	test.NoErr(err)

	results, ok := bodyJson["results"]
	test.True(ok)
	test.True(len(results.([]interface{})) == 1)

	result := results.([]interface{})[0]
	name, ok := result.(map[string]interface{})["name"]
	test.True(ok)
	test.Equal(name, testingCADomain)

	response, err = app.Test(
		httptest.NewRequest("DELETE", "/"+testingCADomain, nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)
}

func TestServerEnumEmpty(t *testing.T) {
	var (
		app  = testApp(t)
		test = is.New(t)
	)

	response, err := app.Test(
		httptest.NewRequest("GET", "/enum", nil),
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

	bodyJson := make(map[string]interface{})
	err = json.Unmarshal(body, &bodyJson)
	test.NoErr(err)

	results, ok := bodyJson["results"]
	test.True(ok && len(results.([]interface{})) == 0)
}
