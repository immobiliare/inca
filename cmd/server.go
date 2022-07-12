package cmd

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
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
				logrus.WithError(err).Fatalln()
			}

			cfg, err := cmd.Flags().GetString("config")
			if err != nil {
				logrus.WithError(err).Fatalln()
			}

			logrus.WithField("bind", bind).Println("Spinning up Inca server")
			inca, err := server.Spinup(cfg)
			if err != nil {
				logrus.WithError(err).Fatalln()
			}

			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			go func() {
				for range sigint {
					logrus.WithError(inca.Shutdown())
				}
			}()

			logrus.WithField("size", len(inca.Cfg.Providers)).Println("Registered providers")
			if err := inca.Listen(bind); err != nil {
				logrus.WithError(err).Fatalln()
			}
		},
	}
)

func init() {
	cmdServer.Flags().StringP("bind", "b", ":8080", "Bind server to interface")
	cmdServer.Flags().StringP("config", "c", "/etc/inca", "Configuration file")
	cmdRoot.AddCommand(cmdServer)
}
