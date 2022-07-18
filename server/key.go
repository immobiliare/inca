package server

import (
	"bytes"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (inca *Inca) handlerKey(c *fiber.Ctx) error {
	_, data, err := (*inca.Cfg.Storage).Get(c.Params("name"))
	if err == nil {
		return c.SendStream(bytes.NewReader(data), len(data))
	}

	log.Info().Str("name", c.Params("name")).Err(err).Msg("cached key not found")
	return c.SendStatus(fiber.StatusNotFound)
}
