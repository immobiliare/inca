package server

import (
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/provider"
	"gitlab.rete.farm/sistemi/inca/server/config"
	"gitlab.rete.farm/sistemi/inca/server/middleware"
	"gitlab.rete.farm/sistemi/inca/storage"
)

type Inca struct {
	*fiber.App
	Storage   *storage.Storage
	Providers []*provider.Provider
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

	if len(cfg.Sentry) > 0 {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Sentry,
			TracesSampleRate: 1.0,
		}); err != nil {
			return nil, fmt.Errorf("sentry: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
		log.Info().Msg("sentry correctly initialized")
	}

	inca := &Inca{
		fiber.New(
			fiber.Config{DisableStartupMessage: true},
		),
		cfg.Storage,
		cfg.Providers,
	}
	inca.Use(compress.New())
	inca.Use(middleware.Logger(zerolog.New(os.Stdout), func(c *fiber.Ctx) bool {
		return c.Path() == "/health"
	}))
	inca.Get("/", inca.handlerEnum)
	inca.Get("/health", inca.handlerHealth)
	inca.Get("/ca/:filter", inca.handlerCA)
	inca.Get("/:name", inca.handlerCRT)
	inca.Get("/:name/key", inca.handlerKey)
	inca.Get("/:name/show", inca.handlerShow)
	inca.Delete("/:name", inca.handlerRevoke)
	return inca, nil
}
