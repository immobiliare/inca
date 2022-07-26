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
	var (
		filter  = c.Params("filter")
		matches = provider.GetByID(filter, inca.Providers)
		match   *provider.Provider
	)
	if matches == nil {
		match = provider.GetByTargetName(filter, nil, inca.Providers)
		if match == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
	} else if len(matches) != 1 {
		log.Error().Str("id", filter).Msg("multiple providers found")
		return c.SendStatus(fiber.StatusBadRequest)
	} else {
		match = matches[0]
	}

	caCrt, err := (*match).CA()
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
