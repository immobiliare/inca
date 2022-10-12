package server

import (
	"archive/zip"
	"bytes"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (inca *Inca) handlerWebDownload(c *fiber.Ctx) error {
	name := c.Params("name")

	if !inca.authorizedTarget(name, c) {
		_ = c.Bind(fiber.Map{"error": "Unauthorized to download the certificate"})
		return inca.handlerWebIndex(c)
	}

	crt, key, err := (*inca.Storage).Get(name)
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Certificate not found"})
		return inca.handlerWebIndex(c)
	}

	out := new(bytes.Buffer)
	zip := zip.NewWriter(out)
	for key, value := range map[string][]byte{
		name + ".crt": crt,
		name + ".key": key,
	} {
		file, err := zip.Create(key)
		if err != nil {
			_ = c.Bind(fiber.Map{"error": "Unable to create ZIP archive entry"})
			log.Error().Err(err).Msg("unable to create ZIP archive entry")
			return inca.handlerWebView(c)
		}

		if _, err := file.Write(value); err != nil {
			_ = c.Bind(fiber.Map{"error": "Unable to add content to ZIP archive entry"})
			log.Error().Err(err).Msg("unable to add content to ZIP archive entry")
			return inca.handlerWebView(c)
		}
	}

	if err := zip.Close(); err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to close ZIP archive"})
		log.Error().Err(err).Msg("unable to close ZIP archive")
		return inca.handlerWebView(c)
	}

	c.Response().Header.Add(fiber.HeaderContentDisposition, `attachment; filename="`+name+`.zip"`)
	return c.SendStream(out, len(out.Bytes()))
}
