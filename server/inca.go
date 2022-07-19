package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gitlab.rete.farm/sistemi/inca/server/config"
	"gitlab.rete.farm/sistemi/inca/server/middleware"
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

	inca := &Inca{fiber.New(fiber.Config{DisableStartupMessage: true}), cfg}
	inca.Use(middleware.Logger(zerolog.New(os.Stdout), func(c *fiber.Ctx) bool { return false }))
	inca.Get("/:name", inca.handlerCRT)
	inca.Get("/:name/key", inca.handlerKey)
	inca.Get("/:name/show", inca.handlerShow)
	inca.Put("/:name/revoke", inca.handlerRevoke)
	inca.Get("/ca/:provider", inca.handlerCA)
	return inca, nil
}
