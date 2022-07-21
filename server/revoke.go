package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (inca *Inca) handlerRevoke(c *fiber.Ctx) error {
	if err := (*inca.Storage).Del(c.Params("name")); err != nil {
		log.Error().Err(err).Msg("unable to remove")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
