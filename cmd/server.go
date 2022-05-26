package cmd

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.rete.farm/sistemi/inca/server"
)

var (
	inca      *server.Inca
	cmdServer = &cobra.Command{
		Use:   "server",
		Short: "Run an Inca server instance",
		Run: func(cmd *cobra.Command, args []string) {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			go func() {
				for range sigint {
					logrus.WithError(inca.Shutdown())
				}
			}()

			bind, err := cmd.Flags().GetString("bind")
			if err != nil {
				logrus.WithError(err).Fatalln()
			}

			inca = server.Spinup()
			if err := inca.Listen(bind); err != nil {
				logrus.WithError(err).Fatalln()
			}
		},
	}
)

func init() {
	cmdServer.Flags().StringP("bind", "b", ":8080", "Bind server to interface")
	cmdRoot.AddCommand(cmdServer)
}
