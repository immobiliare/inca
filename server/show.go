package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/pki"
)

type Certificate struct {
}

func (inca *Inca) handlerShow(c *fiber.Ctx) error {
	var name = c.Params("name")
	if !pki.IsValidCN(name) {
		log.Error().Str("name", name).Msg("invalid name")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	data, _, err := (*inca.Storage).Get(name)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	crt, err := pki.ParseBytes(data)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse certificate")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(EncodeCrt(crt))
}
