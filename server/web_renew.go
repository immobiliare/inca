package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/provider"
	"github.com/rs/zerolog/log"
)

func (inca *Inca) handlerWebRenew(c *fiber.Ctx) error {
	name := c.Params("name")

	if !inca.authorizedTarget(name, c) {
		_ = c.Bind(fiber.Map{"error": "Unauthorized to renew the certificate"})
		return inca.handlerWebIndex(c)
	}

	// Check if certificate exists
	_, _, err := (*inca.Storage).Get(name)
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Certificate not found"})
		return inca.handlerWebIndex(c)
	}

	// Find appropriate provider
	p := provider.GetByTargetName(name, map[string]string{}, inca.Providers)
	if p == nil {
		_ = c.Bind(fiber.Map{"error": "Unable to find a suitable provider for renewal"})
		return inca.handlerWebIndex(c)
	}

	// Get new certificate from provider
	crt, key, err := (*p).Get(name, map[string]string{})
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("unable to renew certificate")
		_ = c.Bind(fiber.Map{"error": "Unable to renew the certificate"})
		return inca.handlerWebIndex(c)
	}

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

	// Check if certificate exists
	_, _, err := (*inca.Storage).Get(name)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	// Find appropriate provider
	p := provider.GetByTargetName(name, map[string]string{}, inca.Providers)
	if p == nil {
		log.Warn().Str("name", name).Msg("no provider found for renewal")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Get new certificate from provider
	crt, key, err := (*p).Get(name, map[string]string{})
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("unable to renew certificate")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Try to renew certificate in place
	if err := (*inca.Storage).Renew(name, crt, key); err != nil {
		log.Error().Err(err).Str("name", name).Msg("unable to renew certificate in storage")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}