package server

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/redirect/v2"
	"github.com/gofiber/template/django"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	templateEngine := django.New(cfg.TemplatesPath, ".html.j2")
	templateEngine.Reload(strings.EqualFold(cfg.Environment, "development"))
	templateEngine.Debug(strings.EqualFold(cfg.Environment, "development"))

	inca := &Inca{
		fiber.New(
			fiber.Config{
				DisableStartupMessage: true,
				Views:                 templateEngine,
				// Views:                 html.NewFileSystem(http.Dir("./server/views"), ".html.j2"),
			},
		),
		cfg.Storage,
		cfg.Providers,
	}
	inca.Use(compress.New())
	inca.Use(middleware.Logger(zerolog.New(os.Stdout), func(c *fiber.Ctx) bool {
		return strings.HasPrefix(c.Path(), "/health") ||
			strings.HasPrefix(c.Path(), "/static/") ||
			strings.HasPrefix(c.Path(), "/favicon.ico")
	}))
	inca.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"^/web$":         "/",
			"^/favicon.ico$": "/static/favicon.ico",
		},
		StatusCode: 301,
	}))

	inca.Get("/", inca.handlerWebIndex)
	inca.Get("/web/config", inca.handlerWebConfig)
	inca.Get("/web/issue", inca.handlerWebIssueView)
	inca.Post("/web/issue", inca.handlerWebIssue)
	inca.Get("/web/import", inca.handlerWebImportView)
	inca.Post("/web/import", inca.handlerWebImport)
	inca.Get("/web/:name", inca.handlerWebView)
	inca.Post("/web/:name/delete", inca.handlerWebDelete)
	inca.Get("/enum", inca.handlerEnum)
	inca.Get("/health", inca.handlerHealth)
	inca.Get("/ca/:filter", inca.handlerCA)
	inca.Get("/:name", inca.handlerCRT)
	inca.Get("/:name/key", inca.handlerKey)
	inca.Get("/:name/show", inca.handlerShow)
	inca.Delete("/:name", inca.handlerRevoke)
	inca.Static("/static", "./server/static")
	return inca, nil
}
