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

func (inca *Inca) handlerCA(c *fiber.Ctx) error {
	p := provider.Get(c.Params("provider"), inca.Cfg.Providers)
	if p == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	caCrt, err := (*p).CA()
	if err != nil {
		log.Error().Err(err).Msg("unable to retrieve CA certificate")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	caCrtBytes := pki.ExportBytes(caCrt)
	if strings.EqualFold(c.Get("Accept", "text/plain"), "application/json") {
		return c.SendString(fmt.Sprintf(`{"crt":"%s"}`, util.BytesToJSON(caCrtBytes)))
	}
	return c.SendStream(bytes.NewReader(caCrtBytes), len(caCrtBytes))
}
