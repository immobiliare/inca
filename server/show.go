package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/pki"
)

func (inca *Inca) handlerShow(c *fiber.Ctx) error {
	var name = c.Params("name")
	if !pki.IsValidCN(name) {
		log.Error().Str("name", name).Msg("invalid name")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	data, _, err := (*inca.Cfg.Storage).Get(name)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	crt, err := pki.ParseBytes(data)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse certificate")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(struct {
		Names     []string
		NotBefore time.Time
		NotAfter  time.Time
	}{
		crt.DNSNames,
		crt.NotBefore,
		crt.NotAfter,
	})
}
