package server

import (
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (inca *Inca) handlerKey(c *fiber.Ctx) error {
	keyFname := keyFilename(c.Params("name"))
	data, err := (*inca.Cfg.Storage).Get(keyFname)
	if err == nil {
		return c.SendStream(bytes.NewReader(data), len(data))
	}

	log.Info().Str("fname", keyFname).Err(err).Msg("cached key not found")
	return c.SendStatus(fiber.StatusNotFound)
}

func keyFilename(name string) string {
	return fmt.Sprintf("%s.key", name)
}
