package server

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/provider"
	"gitlab.rete.farm/sistemi/inca/util"
)

func (inca *Inca) handlerCRT(c *fiber.Ctx) error {
	var (
		name         = c.Params("name")
		queryStrings = util.ParseQueryString(c.Request().URI().QueryString())
	)
	if !pki.IsValidCN(name) {
		log.Error().Str("name", name).Msg("invalid name")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	data, _, err := (*inca.Cfg.Storage).Get(name)
	if err == nil {
		log.Info().Str("name", name).Msg("returning cached certificate")
		if strings.EqualFold(c.Get("Accept", "text/plain"), "application/json") {
			return c.SendString(fmt.Sprintf(`{"crt":"%s"}`, util.BytesToJSON(data)))
		}
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

	if err := (*inca.Cfg.Storage).Put(name, crt, key); err != nil {
		log.Error().Err(err).Msg("unable to persist certificate")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	crtData := pki.ExportBytes(crt)
	if strings.EqualFold(c.Get("Accept", "text/plain"), "application/json") {
		return c.SendString(fmt.Sprintf(`{"crt":"%s"}`, util.BytesToJSON(crtData)))
	}
	return c.SendStream(bytes.NewReader(crtData), len(crtData))
}
