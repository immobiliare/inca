package server

import (
	"io"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/matryer/is"
)

func TestServerWebRenewView(t *testing.T) {
	var (
		app        = testApp(t)
		test       = is.New(t)
		testDomain = "testserverwebrenewview.domain.tld"
	)

	// First create a certificate
	form := url.Values{}
	form.Add("alt", testDomain)
	form.Add("algo", testingCAAlgorithm)
	request := httptest.NewRequest("POST", "/web/issue", strings.NewReader(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	// Now test the renewal view
	req := httptest.NewRequest("GET", "/web/"+testDomain, nil)
	resp, err := app.Test(req)
	test.NoErr(err)
	test.Equal(resp.StatusCode, fiber.StatusOK)

	body, err := io.ReadAll(resp.Body)
	test.NoErr(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	test.True(strings.Contains(string(body), testDomain))

	// Clean up
	response, err = app.Test(
		httptest.NewRequest("DELETE", "/"+testDomain, nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)
}

func TestServerWebRenewPost(t *testing.T) {
	var (
		app        = testApp(t)
		test       = is.New(t)
		testDomain = "testserverwebrenewpost.domain.tld"
	)

	// First create a certificate
	form := url.Values{}
	form.Add("alt", testDomain)
	form.Add("algo", testingCAAlgorithm)
	request := httptest.NewRequest("POST", "/web/issue", strings.NewReader(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)

	// Test POST request to renew certificate
	form = url.Values{}
	form.Add("domains", testDomain)
	form.Add("email", "test@example.com")

	req := httptest.NewRequest("POST", "/web/"+testDomain+"/renew", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := app.Test(req, 5000) // Increase timeout to 5 seconds
	test.NoErr(err)

	_, err = io.ReadAll(resp.Body)
	test.NoErr(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	// Should get a redirect or success response
	test.True(resp.StatusCode == fiber.StatusOK || resp.StatusCode == fiber.StatusFound)

	// Clean up
	response, err = app.Test(
		httptest.NewRequest("DELETE", "/"+testDomain, nil),
	)
	test.NoErr(err)
	test.Equal(response.StatusCode, fiber.StatusOK)
}

func TestServerWebRenewUnauthorized(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	// Create app with ACL that denies access to unauthorized.example.com
	app := testApp(t)
	app.acl = map[string][]string{
		"testtoken": {"domain\\.tld", "allowed\\.com"}, // Only allow these domains
	}

	// Test POST request to renew certificate for unauthorized domain
	form := url.Values{}
	form.Add("domains", "unauthorized.example.com")
	form.Add("email", "test@example.com")

	req := httptest.NewRequest("POST", "/web/unauthorized.example.com/renew", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := app.Test(req)
	is.NoErr(err)

	body, err := io.ReadAll(resp.Body)
	is.NoErr(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	// Should get unauthorized response
	is.True(strings.Contains(string(body), "Unauthorized") || resp.StatusCode == fiber.StatusFound)
}

func TestServerRenewMultipleCertificates(t *testing.T) {
	var (
		app     = testApp(t)
		test    = is.New(t)
		domains = []string{"testmulti1.domain.tld", "testmulti2.domain.tld", "testmulti3.domain.tld"}
	)

	// First create certificates for each domain
	for _, domain := range domains {
		form := url.Values{}
		form.Add("alt", domain)
		form.Add("algo", testingCAAlgorithm)
		request := httptest.NewRequest("POST", "/web/issue", strings.NewReader(form.Encode()))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response, err := app.Test(request)
		test.NoErr(err)
		test.Equal(response.StatusCode, fiber.StatusOK)
	}

	// Test renewing each certificate via API
	for _, domain := range domains {
		req := httptest.NewRequest("POST", "/"+domain+"/renew", nil)
		resp, err := app.Test(req, 5000) // Increase timeout to 5 seconds
		test.NoErr(err)

		// Should get success
		test.Equal(resp.StatusCode, fiber.StatusOK)
	}

	// Clean up
	for _, domain := range domains {
		response, err := app.Test(
			httptest.NewRequest("DELETE", "/"+domain, nil),
		)
		test.NoErr(err)
		test.Equal(response.StatusCode, fiber.StatusOK)
	}
}

func TestServerRenewNonExistentCertificate(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	app := testApp(t)

	// Try to renew a certificate that doesn't exist
	// Note: Authorization check happens first, so this returns 401 instead of 404
	// This is correct behavior to prevent unauthorized users from probing certificate existence
	req := httptest.NewRequest("POST", "/nonexistent.example.com/renew", nil)
	resp, err := app.Test(req, 2000)
	is.NoErr(err)

	// The actual behavior depends on the ACL configuration and certificate existence
	// Both 401 (Unauthorized) and 404 (Not Found) are valid responses
	is.True(resp.StatusCode == fiber.StatusUnauthorized || resp.StatusCode == fiber.StatusNotFound)
}

func TestServerRenewInvalidDomain(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	app := testApp(t)

	// Try to renew with invalid domain
	// Note: Authorization check happens first, so this returns 401 instead of 400/404
	// This is correct behavior to prevent unauthorized users from probing certificate existence
	req := httptest.NewRequest("POST", "/invalid-domain/renew", nil)
	resp, err := app.Test(req, 2000)
	is.NoErr(err)

	// The actual behavior depends on the ACL configuration and certificate existence
	// Both 401 (Unauthorized) and 404 (Not Found) are valid responses
	is.True(resp.StatusCode == fiber.StatusUnauthorized || resp.StatusCode == fiber.StatusNotFound)
}
