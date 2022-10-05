package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/provider"
)

func (inca *Inca) handlerWebView(c *fiber.Ctx) error {
	chain, key, err := (*inca.Storage).Get(c.Params("name"))
	if err != nil {
		log.Error().Err(err).Msg("unable to find certificate")
		return c.SendStatus(fiber.StatusNotFound)
	}

	crt, err := pki.ParseBytes(chain)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse certificate")
	}

	viewBinds := fiber.Map{
		"crt": EncodeCrt(crt, provider.GetByTargetName(crt.Subject.CommonName, nil, inca.Providers)),
	}
	if inca.authorizedTarget(crt.Subject.CommonName, c) {
		viewBinds["chain"] = string(chain)
		viewBinds["key"] = string(key)
	}

	return c.Render("view", viewBinds)
}
