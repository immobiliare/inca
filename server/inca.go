package server

import (
	"crypto/x509"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/server/config"
	"gitlab.rete.farm/sistemi/inca/server/middleware"
)

type Inca struct {
	*fiber.App
	Cfg *config.Config
}

type Crt struct {
	CN        string    `json:"name"`
	AltNames  []string  `json:"alt"`
	NotBefore time.Time `json:"not_before"`
	NotAfter  time.Time `json:"not_after"`
}

func EncodeCrt(crt *x509.Certificate) Crt {
	dnsNames, ipAddresses := pki.AltNames(crt)
	return Crt{
		crt.Subject.CommonName,
		append(dnsNames, ipAddresses...),
		crt.NotBefore,
		crt.NotAfter,
	}
}

func Spinup(path string) (*Inca, error) {
	cfg, err := config.Parse(path)
	if err != nil {
		return nil, err
	}

	inca := &Inca{fiber.New(fiber.Config{DisableStartupMessage: true}), cfg}
	inca.Use(middleware.Logger(zerolog.New(os.Stdout), func(c *fiber.Ctx) bool { return false }))
	inca.Get("/", inca.handlerEnum)
	inca.Get("/:name", inca.handlerCRT)
	inca.Get("/:name/key", inca.handlerKey)
	inca.Get("/:name/show", inca.handlerShow)
	inca.Put("/:name/revoke", inca.handlerRevoke)
	inca.Get("/ca/:provider", inca.handlerCA)
	return inca, nil
}
