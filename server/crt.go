package server

import (
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/provider"
	"gitlab.rete.farm/sistemi/inca/util"
)

func (inca *Inca) handlerCRT(c *fiber.Ctx) error {
	var (
		name         = c.Params("name")
		crtFname     = crtFilename(name)
		keyFname     = keyFilename(name)
		queryStrings = util.ParseQueryString(c.Request().URI().QueryString())
	)
	if len(name) <= 3 {
		log.Error().Str("name", name).Msg("name too short")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	data, err := (*inca.Cfg.Storage).Get(crtFname)
	if err == nil {
		log.Info().Str("fname", crtFname).Err(err).Msg("returning cached certificate")
		return c.SendStream(bytes.NewReader(data), len(data))
	}

	p := provider.GetFor(name, queryStrings, (*inca.Cfg).Providers)
	if p == nil {
		log.Error().Str("name", name).Msg("no provider found")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	crt, key, err := (*p).Get(name, queryStrings)
	if err != nil {
		log.Error().Err(err).Msg("unable to generate")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := (*inca.Cfg.Storage).Put(crtFname, crt); err != nil {
		log.Error().Err(err).Msg("unable to persist certificate")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := (*inca.Cfg.Storage).Put(keyFname, key); err != nil {
		log.Error().Err(err).Msg("unable to persist certificate")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStream(bytes.NewReader(crt.Bytes), len(crt.Bytes))
}

func crtFilename(name string) string {
	return fmt.Sprintf("%s.pem", name)
}
