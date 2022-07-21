package cmd

import (
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gitlab.rete.farm/sistemi/inca/server"
)

var (
	cmdServer = &cobra.Command{
		Use:   "server",
		Short: "Run an Inca server instance",
		Run: func(cmd *cobra.Command, args []string) {
			bind, err := cmd.Flags().GetString("bind")
			if err != nil {
				log.Fatal().Err(err).Msg("bind flag is mandatory")
			}

			cfg, err := cmd.Flags().GetString("config")
			if err != nil {
				log.Fatal().Err(err).Msg("config flag is mandatory")
			}

			log.Info().Str("bind", bind).Msg("spinning up Inca server")
			inca, err := server.Spinup(cfg)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to spinup Inca server")
			}

			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			go func() {
				for range sigint {
					log.Error().Err(inca.Shutdown())
				}
			}()

			log.Info().Int("size", len(inca.Providers)).Msg("registered providers")
			log.Info().Str("type", (*(inca.Storage)).ID()).Msg("registered storage")
			if err := inca.Listen(bind); err != nil {
				log.Fatal().Err(err).Msg("unable to bind server")
			}
		},
	}
)

func init() {
	cmdServer.Flags().StringP("bind", "b", ":8080", "Bind server to interface")
	cmdServer.Flags().StringP("config", "c", "/etc/inca", "Configuration file")
	cmdRoot.AddCommand(cmdServer)
}
