package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (inca *Inca) handlerRevoke(c *fiber.Ctx) error {
	for _, asset := range []string{
		crtFilename(c.Params("name")),
		keyFilename(c.Params("name")),
	} {
		if err := (*inca.Cfg.Storage).Del(asset); err != nil {
			log.Error().Err(err).Msg("unable to remove")
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
