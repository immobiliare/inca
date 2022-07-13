package server

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/provider"
	"gitlab.rete.farm/sistemi/inca/server/config"
	"gitlab.rete.farm/sistemi/inca/util"
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

	inca := &Inca{
		fiber.New(
			fiber.Config{DisableStartupMessage: true},
		),
		cfg,
	}
	inca.Use(zlogger(zerolog.New(os.Stdout), func(c *fiber.Ctx) bool { return false }))
	inca.Get("/:name", func(c *fiber.Ctx) error {
		var (
			name         = c.Params("name")
			queryStrings = util.ParseQueryString(c.Request().URI().QueryString())
		)
		if len(name) <= 3 {
			log.Error().Str("name", name).Msg("name too short")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		p := provider.GetFor(name, queryStrings, cfg.Providers)
		if p == nil {
			log.Error().Str("name", name).Msg("no provider found")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		crt, key, err := (*p).Get(name, queryStrings)
		if err != nil {
			log.Error().Err(err).Msg("unable to generate")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := pki.Export(crt, fmt.Sprintf("%s.pem", name), false); err != nil {
			log.Error().Err(err).Msg("unable to export certificate")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := pki.Export(key, fmt.Sprintf("%s.key", name), true); err != nil {
			log.Error().Err(err).Msg("unable to export key")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.SendString("ok")
	})
	inca.Get("/revoke/:name", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("revoke %s", c.Params("name")))
	})

	return inca, nil
}
