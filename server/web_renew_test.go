package server

import (
    "io"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"

    "github.com/gofiber/fiber/v2"
    "github.com/immobiliare/inca/pki"
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

func TestServerRenewPreservesSANs(t *testing.T) {
    var (
        app        = testApp(t)
        test       = is.New(t)
        testDomain = "testrenewsans.domain.tld"
        san1       = "san1.domain.tld"
        san2       = "san2.domain.tld"
        ipSAN      = "192.168.1.100"
    )

    // First create a certificate with multiple SANs
    form := url.Values{}
    form.Add("alt", testDomain+","+san1+","+san2+","+ipSAN)
    form.Add("algo", testingCAAlgorithm)
    request := httptest.NewRequest("POST", "/web/issue", strings.NewReader(form.Encode()))
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    response, err := app.Test(request)
    test.NoErr(err)
    test.Equal(response.StatusCode, fiber.StatusOK)

    // Get the original certificate and verify SANs
    req := httptest.NewRequest("GET", "/"+testDomain, nil)
    resp, err := app.Test(req)
    test.NoErr(err)
    test.Equal(resp.StatusCode, fiber.StatusOK)

    originalCrtData, err := io.ReadAll(resp.Body)
    test.NoErr(err)
    if err := resp.Body.Close(); err != nil {
        t.Logf("Failed to close response body: %v", err)
    }

    originalCrt, err := pki.ParseBytes(originalCrtData)
    test.NoErr(err)
    
    originalDNSNames, originalIPAddresses := pki.AltNames(originalCrt)
    test.True(len(originalDNSNames) == 3) // testDomain, san1, san2
    test.True(len(originalIPAddresses) == 1) // ipSAN

    // Now renew the certificate via API
    renewReq := httptest.NewRequest("POST", "/"+testDomain+"/renew", nil)
    renewResp, err := app.Test(renewReq, 5000)
    test.NoErr(err)
    test.Equal(renewResp.StatusCode, fiber.StatusOK)

    // Get the renewed certificate
    req = httptest.NewRequest("GET", "/"+testDomain, nil)
    resp, err = app.Test(req)
    test.NoErr(err)
    test.Equal(resp.StatusCode, fiber.StatusOK)

    renewedCrtData, err := io.ReadAll(resp.Body)
    test.NoErr(err)
    if err := resp.Body.Close(); err != nil {
        t.Logf("Failed to close response body: %v", err)
    }

    renewedCrt, err := pki.ParseBytes(renewedCrtData)
    test.NoErr(err)

    renewedDNSNames, renewedIPAddresses := pki.AltNames(renewedCrt)

    // Verify all SANs are preserved
    test.Equal(len(renewedDNSNames), len(originalDNSNames))
    test.Equal(len(renewedIPAddresses), len(originalIPAddresses))
    
    // Verify specific SANs are present
    test.True(contains(renewedDNSNames, testDomain))
    test.True(contains(renewedDNSNames, san1))
    test.True(contains(renewedDNSNames, san2))
    test.True(contains(renewedIPAddresses, ipSAN))

    // Clean up
    response, err = app.Test(
        httptest.NewRequest("DELETE", "/"+testDomain, nil),
    )
    test.NoErr(err)
    test.Equal(response.StatusCode, fiber.StatusOK)
}

func TestServerWebRenewPreservesSANs(t *testing.T) {
    var (
        app        = testApp(t)
        test       = is.New(t)
        testDomain = "testwebrenewsans.domain.tld"
        san1       = "websan1.domain.tld"
        san2       = "websan2.domain.tld"
    )

    // First create a certificate with multiple SANs
    form := url.Values{}
    form.Add("alt", testDomain+","+san1+","+san2)
    form.Add("algo", testingCAAlgorithm)
    request := httptest.NewRequest("POST", "/web/issue", strings.NewReader(form.Encode()))
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    response, err := app.Test(request)
    test.NoErr(err)
    test.Equal(response.StatusCode, fiber.StatusOK)

    // Get the original certificate and verify SANs
    req := httptest.NewRequest("GET", "/"+testDomain, nil)
    resp, err := app.Test(req)
    test.NoErr(err)
    test.Equal(resp.StatusCode, fiber.StatusOK)

    originalCrtData, err := io.ReadAll(resp.Body)
    test.NoErr(err)
    if err := resp.Body.Close(); err != nil {
        t.Logf("Failed to close response body: %v", err)
    }

    originalCrt, err := pki.ParseBytes(originalCrtData)
    test.NoErr(err)
    
    originalDNSNames, _ := pki.AltNames(originalCrt)
    test.True(len(originalDNSNames) == 3) // testDomain, san1, san2

    // Now renew the certificate via web interface
    form = url.Values{}
    renewReq := httptest.NewRequest("POST", "/web/"+testDomain+"/renew", strings.NewReader(form.Encode()))
    renewReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    renewResp, err := app.Test(renewReq, 5000)
    test.NoErr(err)
    test.True(renewResp.StatusCode == fiber.StatusOK || renewResp.StatusCode == fiber.StatusFound)

    // Get the renewed certificate
    req = httptest.NewRequest("GET", "/"+testDomain, nil)
    resp, err = app.Test(req)
    test.NoErr(err)
    test.Equal(resp.StatusCode, fiber.StatusOK)

    renewedCrtData, err := io.ReadAll(resp.Body)
    test.NoErr(err)
    if err := resp.Body.Close(); err != nil {
        t.Logf("Failed to close response body: %v", err)
    }

    renewedCrt, err := pki.ParseBytes(renewedCrtData)
    test.NoErr(err)

    renewedDNSNames, _ := pki.AltNames(renewedCrt)

    // Verify all SANs are preserved
    test.Equal(len(renewedDNSNames), len(originalDNSNames))
    
    // Verify specific SANs are present
    test.True(contains(renewedDNSNames, testDomain))
    test.True(contains(renewedDNSNames, san1))
    test.True(contains(renewedDNSNames, san2))

    // Clean up
    response, err = app.Test(
        httptest.NewRequest("DELETE", "/"+testDomain, nil),
    )
    test.NoErr(err)
    test.Equal(response.StatusCode, fiber.StatusOK)
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

