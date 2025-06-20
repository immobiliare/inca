package server

import (
	"archive/zip"
	"bytes"
	"crypto/x509"
	"encoding/pem"

	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/util"
	"github.com/rs/zerolog/log"
	"software.sslmate.com/src/go-pkcs12"
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

func (inca *Inca) handlerWebDownloadPfx(c *fiber.Ctx) error {
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

	certBlock, _ := pem.Decode([]byte(crt))
	if certBlock == nil {
		_ = c.Bind(fiber.Map{"error": "Invalid PEM content for certificate"})
		log.Error().Msg("invalid PEM content for certificate")
		return inca.handlerWebView(c)
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to parse certificate"})
		log.Error().Err(err).Msg("unable to parse certificate")
		return inca.handlerWebView(c)
	}

	keyBlock, _ := pem.Decode([]byte(key))
	privateKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to parse private key"})
		log.Error().Err(err).Msg("unable to parse private key")
		return inca.handlerWebView(c)
	}

	password := util.GenerateRandomString(256)
	pfxData, err := pkcs12.Modern.Encode(privateKey, cert, nil, password)
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to create PFX"})
		log.Error().Err(err).Msg("unable to create PFX")
		return inca.handlerWebView(c)
	}

	out := new(bytes.Buffer)
	zip := zip.NewWriter(out)
	for key, value := range map[string][]byte{
		name + ".pfx": pfxData,
		name + ".txt": []byte(password),
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
