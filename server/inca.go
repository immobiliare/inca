package server

import (
	"bytes"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/provider"
	"gitlab.rete.farm/sistemi/inca/server/config"
	"gitlab.rete.farm/sistemi/inca/server/middleware"
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
	inca.Use(middleware.Logger(zerolog.New(os.Stdout), func(c *fiber.Ctx) bool { return false }))
	inca.Get("/:name", func(c *fiber.Ctx) error {
		var (
			name         = c.Params("name")
			crtFname     = fmt.Sprintf("%s.pem", name)
			keyFname     = fmt.Sprintf("%s.key", name)
			queryStrings = util.ParseQueryString(c.Request().URI().QueryString())
		)
		if len(name) <= 3 {
			log.Error().Str("name", name).Msg("name too short")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		data, err := (*cfg.Storage).Get(crtFname)
		if err == nil {
			log.Info().Str("fname", crtFname).Err(err).Msg("returning cached certificate")
			return c.SendStream(bytes.NewReader(data), len(data))
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

		if err := (*cfg.Storage).Put(crtFname, crt); err != nil {
			log.Error().Err(err).Msg("unable to persist certificate")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := (*cfg.Storage).Put(keyFname, key); err != nil {
			log.Error().Err(err).Msg("unable to persist certificate")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.SendStream(bytes.NewReader(crt.Bytes), len(crt.Bytes))
	})
	inca.Get("/:name/key", func(c *fiber.Ctx) error {
		var keyFname = fmt.Sprintf("%s.key", c.Params("name"))
		data, err := (*cfg.Storage).Get(keyFname)
		if err == nil {
			return c.SendStream(bytes.NewReader(data), len(data))
		}

		log.Info().Str("fname", keyFname).Err(err).Msg("cached key not found")
		return c.SendStatus(fiber.StatusNotFound)
	})
	inca.Get("/ca/:provider", func(c *fiber.Ctx) error {
		p := provider.Get(c.Params("provider"), inca.Cfg.Providers)
		if p == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		caCrt, err := (*p).CA()
		if err != nil {
			log.Error().Err(err).Msg("unable to retrieve CA certificate")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.SendStream(bytes.NewReader(caCrt.Bytes), len(caCrt.Bytes))
	})
	inca.Get("/revoke/:name", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("revoke %s", c.Params("name")))
	})

	return inca, nil
}
