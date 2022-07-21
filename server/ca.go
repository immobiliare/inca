package server

import (
	"bytes"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/provider"
)

func (inca *Inca) handlerCA(c *fiber.Ctx) error {
	p := provider.GetFrom(c.Params("provider"), inca.Providers)
	if p == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	caCrt, err := (*p).CA()
	if err != nil {
		log.Error().Err(err).Msg("unable to retrieve CA certificate")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	caCrtBytes := pki.ExportBytes(caCrt)
	if strings.EqualFold(c.Get("Accept", "text/plain"), "application/json") {
		return c.JSON(struct {
			Crt string `json:"crt"`
		}{string(caCrtBytes)})
	}
	return c.SendStream(bytes.NewReader(caCrtBytes), len(caCrtBytes))
}
