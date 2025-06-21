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

	var certs []*x509.Certificate
	rest := []byte(crt)
	for {
		var certBlock *pem.Block
		certBlock, rest = pem.Decode(rest)
		if certBlock == nil {
			break
		}
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			_ = c.Bind(fiber.Map{"error": "Unable to parse certificate"})
			log.Error().Err(err).Msg("unable to parse certificate")
			return inca.handlerWebView(c)
		}
		certs = append(certs, cert)
	}
	if len(certs) == 0 {
		_ = c.Bind(fiber.Map{"error": "No valid certificates found"})
		log.Error().Msg("no valid certificates found")
		return inca.handlerWebView(c)
	}

	leaf := certs[0]
	for _, candidate := range certs {
		isIssuer := false
		for _, other := range certs {
			if candidate.Subject.String() == other.Issuer.String() && candidate != other {
				isIssuer = true
				break
			}
		}
		if !isIssuer {
			leaf = candidate
			break
		}
	}

	var chain []*x509.Certificate
	for _, cert := range certs {
		if cert != leaf {
			chain = append(chain, cert)
		}
	}

	keyBlock, _ := pem.Decode([]byte(key))
	if keyBlock == nil {
		_ = c.Bind(fiber.Map{"error": "Invalid PEM data for private key"})
		log.Error().Msg("invalid PEM data for private key")
		return inca.handlerWebView(c)
	}

	var privateKey interface{}
	var parseErr error

	privateKey, parseErr = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if parseErr != nil {
		privateKey, parseErr = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if parseErr != nil {
			privateKey, parseErr = x509.ParseECPrivateKey(keyBlock.Bytes)
			if parseErr != nil {
				_ = c.Bind(fiber.Map{"error": "Unable to parse private key"})
				log.Error().Err(parseErr).Msg("unable to parse private key")
				return inca.handlerWebView(c)
			}
		}
	}

	password := util.GenerateRandomString(256)
	pfxData, err := pkcs12.Modern.Encode(privateKey, leaf, chain, password)
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
