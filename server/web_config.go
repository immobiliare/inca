package server

import (
	"github.com/gofiber/fiber/v2"
)

func (inca *Inca) handlerWebConfig(c *fiber.Ctx) error {
	return c.Render("config", fiber.Map{
		"storage":   EncodeStorage(inca.Storage),
		"providers": EncodeProviders(inca.Providers),
	})
}
