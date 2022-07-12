package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gitlab.rete.farm/sistemi/inca/server/config"
)

const (
	logFormat = "{\"rtime\":\"${latency}\",\"ip\":\"${ip}\",\"port\":\"${port}\",\"status\":\"${status}\",\"method\":\"${method}\",\"path\":\"${path}\",\"time\":\"${time}\"}\n"
)

type Inca struct {
	*fiber.App
	Cfg *config.Config
}

func Spinup(path string) (*Inca, error) {
	cfg, err := config.Parse(path)
	if err != nil {
		return nil, err
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(logger.New(logger.Config{
		Format:     logFormat,
		TimeFormat: time.RFC3339,
	}))
	app.Get("/:name", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("get %s", c.Params("name")))
	})
	app.Get("/revoke/:name", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("revoke %s", c.Params("name")))
	})

	return &Inca{app, cfg}, nil
}
