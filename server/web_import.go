package server

import (
	"bytes"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/pki"
)

func (inca *Inca) handlerWebImportView(c *fiber.Ctx) error {
	return c.Render("import", fiber.Map{})
}

func (inca *Inca) handlerWebImport(c *fiber.Ctx) error {
	formCrtHeader, err := c.FormFile("crt")
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "No certificate given"})
		return inca.handlerWebImportView(c)
	}

	formCrtFile, err := formCrtHeader.Open()
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to open certificate file"})
		return inca.handlerWebImportView(c)
	}

	formCrtBuf := bytes.NewBuffer(nil)
	if _, err := io.Copy(formCrtBuf, formCrtFile); err != nil {
		log.Error().Err(err).Msg("unable to read certificate file")
		_ = c.Bind(fiber.Map{"error": "Unable to read certificate file"})
		return inca.handlerWebImportView(c)
	}

	formKeyHeader, err := c.FormFile("key")
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "No key given"})
		return inca.handlerWebImportView(c)
	}

	formKeyFile, err := formKeyHeader.Open()
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to open key file"})
		return inca.handlerWebImportView(c)
	}

	formKeyBuf := bytes.NewBuffer(nil)
	if _, err := io.Copy(formKeyBuf, formKeyFile); err != nil {
		log.Error().Err(err).Msg("unable to read key file")
		_ = c.Bind(fiber.Map{"error": "Unable to read key file"})
		return inca.handlerWebImportView(c)
	}

	crt, _, err := pki.ParseKeyPairBytes(formCrtBuf.Bytes(), formKeyBuf.Bytes())
	if err != nil {
		log.Error().Err(err).Msg("unable to parse key pair")
		_ = c.Bind(fiber.Map{"error": "Unable to parse key pair"})
		return inca.handlerWebImportView(c)
	}

	if _, _, err := (*inca.Storage).Get(crt.Subject.CommonName); err == nil {
		_ = c.Bind(fiber.Map{"error": "Certificate already existing"})
		return inca.handlerWebImportView(c)
	}

	if err := (*inca.Storage).Put(crt.Subject.CommonName, formCrtBuf.Bytes(), formKeyBuf.Bytes()); err != nil {
		log.Error().Err(err).Msg("unable to persist certificate")
		_ = c.Bind(fiber.Map{"error": "Unable to persist the certificate"})
		return inca.handlerWebImportView(c)
	}

	_ = c.Bind(fiber.Map{"message": "Certificate successfully imported: " + crt.Subject.CommonName})
	return inca.handlerWebIndex(c)
}
