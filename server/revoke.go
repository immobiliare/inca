package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/provider"
)

func (inca *Inca) handlerRevoke(c *fiber.Ctx) error {
	name := c.Params("name")
	if !inca.authorizedTarget(name, c) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	data, _, err := (*inca.Storage).Get(name)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if err := (*inca.Storage).Del(name); err != nil {
		log.Error().Err(err).Msg("unable to remove")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	p := provider.GetByTargetName(name, map[string]string{}, inca.Providers)
	if p == nil {
		log.Warn().Str("name", name).Msg("no provider found")
	} else {
		if err := (*p).Del(name, data); err != nil {
			log.Error().Err(err).Str("name", name).Msg("unable to revoke")
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
