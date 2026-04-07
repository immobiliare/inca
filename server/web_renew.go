package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/pki"
	"github.com/immobiliare/inca/provider"
	"github.com/rs/zerolog/log"
)

// CertificateResult holds the result of certificate generation
type CertificateResult struct {
	Crt []byte
	Key []byte
	Err error
}

// getProvider is a shared helper method that handles provider lookup
func (inca *Inca) getProvider(name string, options map[string]string) *provider.Provider {
	return provider.GetByTargetName(name, options, inca.Providers)
}

// generateCertificate is a shared helper method that handles provider lookup and certificate generation
func (inca *Inca) generateCertificate(name string, options map[string]string) CertificateResult {
	// Find appropriate provider
	p := inca.getProvider(name, options)
	if p == nil {
		log.Warn().Str("name", name).Msg("no provider found")
		return CertificateResult{Err: fiber.NewError(fiber.StatusBadRequest, "Unable to find a suitable provider")}
	}

	// Get certificate from provider
	crt, key, err := (*p).Get(name, options)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("unable to generate certificate")
		return CertificateResult{Err: fiber.NewError(fiber.StatusInternalServerError, "Unable to generate certificate")}
	}

	return CertificateResult{Crt: crt, Key: key, Err: nil}
}

func (inca *Inca) handlerWebRenew(c *fiber.Ctx) error {
	name := c.Params("name")

	if !inca.authorizedTarget(name, c) {
		_ = c.Bind(fiber.Map{"error": "Unauthorized to renew the certificate"})
		return inca.handlerWebIndex(c)
	}

	// Check if certificate exists
	crtData, _, err := (*inca.Storage).Get(name)
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Certificate not found"})
		return inca.handlerWebIndex(c)
	}

	// Parse existing certificate to preserve SANs
    options := map[string]string{}
    if crt, err := pki.ParseBytes(crtData); err == nil {
        dnsNames, ipAddresses := pki.AltNames(crt)
        allAltNames := append(dnsNames, ipAddresses...)
        if len(allAltNames) > 0 {
            options["alt"] = strings.Join(allAltNames, ",")
        }
    }

    // Generate new certificate with preserved SANs
    result := inca.generateCertificate(name, options)
    if result.Err != nil {
        _ = c.Bind(fiber.Map{"error": "Unable to renew the certificate"})
        return inca.handlerWebIndex(c)
    }
    crt, key := result.Crt, result.Key

    // Try to renew certificate in place
    if err := (*inca.Storage).Renew(name, crt, key); err != nil {
        log.Error().Err(err).Str("name", name).Msg("unable to renew certificate in storage")
        _ = c.Bind(fiber.Map{"error": "Unable to renew certificate in storage"})
        return inca.handlerWebIndex(c)
    }

    _ = c.Bind(fiber.Map{"message": "Certificate successfully renewed"})
    return inca.handlerWebIndex(c)
}

// API endpoint for certificate renewal
func (inca *Inca) handlerRenew(c *fiber.Ctx) error {
    name := c.Params("name")
    if !inca.authorizedTarget(name, c) {
        return c.SendStatus(fiber.StatusUnauthorized)
    }

    // Check if certificate exists and get its SANs
    crtData, _, err := (*inca.Storage).Get(name)
    if err != nil {
        return c.SendStatus(fiber.StatusNotFound)
    }

    // Parse existing certificate to preserve SANs
    options := map[string]string{}
    if crt, err := pki.ParseBytes(crtData); err == nil {
        dnsNames, ipAddresses := pki.AltNames(crt)
        allAltNames := append(dnsNames, ipAddresses...)
        if len(allAltNames) > 0 {
            options["alt"] = strings.Join(allAltNames, ",")
        }
    }

    // Generate new certificate with preserved SANs
    result := inca.generateCertificate(name, options)
    if result.Err != nil {
        if fiberErr, ok := result.Err.(*fiber.Error); ok {
            return c.SendStatus(fiberErr.Code)
        }
        return c.SendStatus(fiber.StatusInternalServerError)
    }
    crt, key := result.Crt, result.Key

    // Try to renew certificate in place
    if err := (*inca.Storage).Renew(name, crt, key); err != nil {
        log.Error().Err(err).Str("name", name).Msg("unable to renew certificate in storage")
        return c.SendStatus(fiber.StatusInternalServerError)
    }

    return c.SendStatus(fiber.StatusOK)
}