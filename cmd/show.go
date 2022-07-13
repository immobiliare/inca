package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gitlab.rete.farm/sistemi/inca/pki"
)

var showGen = &cobra.Command{
	Use:   "show [certificate] [key]",
	Short: "Pretty print a certificate",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal().Msg("wrong arguments")
		}
		var (
			crtPath = args[0]
			keyPath = args[1]
		)

		log.Info().Str("path", keyPath).Msg("parsing key")
		key, err := pki.ParseKey(keyPath)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to parse key")
		}
		log.Info().Str("path", keyPath).Msg("key parsed")
		log.Info().Str("public", fmt.Sprintf("%v", key.Public())).Msg("key parsed")

		log.Info().Str("path", crtPath).Msg("parsing certificate")
		crt, err := pki.Parse(crtPath)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to parse certificate")
		}
		log.Info().Strs("names", crt.DNSNames).Msg("certificate parsed")
	},
}

func init() {
	cmdRoot.AddCommand(showGen)
}
