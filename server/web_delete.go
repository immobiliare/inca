package server

import (
	"github.com/gofiber/fiber/v2"
)

func (inca *Inca) handlerWebDelete(c *fiber.Ctx) error {
	name := c.Params("name")

	if !inca.authorizedTarget(name, c) {
		_ = c.Bind(fiber.Map{"error": "Unauthorized to delete the certificate"})
		return inca.handlerWebIndex(c)
	}

	data, _, err := (*inca.Storage).Get(name)
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Certificate not found"})
		return inca.handlerWebIndex(c)
	}

	if err := (*inca.Storage).Del(name); err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to remove"})
		return inca.handlerWebIndex(c)
	}

	p := inca.getProvider(name, map[string]string{})
	if p == nil {
		_ = c.Bind(fiber.Map{"error": "Provider not found: certificate not revoked"})
		return inca.handlerWebIndex(c)
	} else {
		if err := (*p).Del(name, data); err != nil {
			_ = c.Bind(fiber.Map{"error": "Unable to revoke"})
			return inca.handlerWebIndex(c)
		}
	}

	_ = c.Bind(fiber.Map{"message": "Certificate successfully deleted"})
	return inca.handlerWebIndex(c)
}
