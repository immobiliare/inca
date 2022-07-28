package server

import (
	"github.com/gofiber/fiber/v2"
)

func (inca *Inca) handlerHealth(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
