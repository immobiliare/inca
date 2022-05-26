package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Inca struct {
	*fiber.App
}

func Spinup() *Inca {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Get("/:name", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("get %s", c.Params("name")))
	})
	app.Get("/revoke/:name", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("revoke %s", c.Params("name")))
	})
	return &Inca{app}
}
