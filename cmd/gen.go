package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/util"
)

var cmdGen = &cobra.Command{
	Use:   "gen",
	Short: "Generate CA certificate",
	Run: func(cmd *cobra.Command, args []string) {
		names, err := cmd.Flags().GetStringArray("name")
		if err != nil {
			log.Fatal().Err(err).Msg("at least a name gotta be given")
		}

		req := pki.NewRequest(names...)
		req.CA = true

		output, err := cmd.Flags().GetString("out")
		if err != nil {
			log.Fatal().Err(err).Msg("output flag is mandatory")
		}

		algo, err := cmd.Flags().GetString("algo")
		if err != nil {
			log.Fatal().Err(err).Msg("algorithm flag is mandatory")
		}
		req.Algo = map[string]int{
			"eddsa": pki.EDDSA,
			"ecdsa": pki.ECDSA,
			"rsa":   pki.RSA,
		}[algo]
		log.Info().Strs("names", req.Hosts).Dur("duration", req.Duration).Str("algo", algo).Msg("generating certificate")
		crt, key, err := pki.New(req)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to generate certificate")
		}
		log.Info().Str("certificate", crt.SerialNumber.String()).Msg("certificate generated")
		log.Info().Str("public", fmt.Sprintf("%v", key.Public())).Msg("key generated")

		log.Info().Msg("wrapping certificate")
		crtBytes, keyBytes, err := pki.Wrap(crt, key, crt, key)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to wrap certificate")
		}

		log.Info().Msg("exporting certificate")
		if err := pki.Export(crtBytes, filepath.Join(output, "crt.pem")); err != nil {
			log.Fatal().Err(err).Msg("unable to export certificate")
		}

		log.Info().Msg("exporting key")
		if err := pki.Export(keyBytes, filepath.Join(output, "key.pem")); err != nil {
			log.Fatal().Err(err).Msg("unable to export key")
		}

		log.Info().Str("output", output).Msg("certificate created")
	},
}

func init() {
	cmdRoot.AddCommand(cmdGen)
	cmdGen.Flags().StringP("out", "o", util.ErrWrap("./")(os.Getwd()), "Output path")
	cmdGen.Flags().StringArrayP("name", "n", []string{}, "Certificate names")
	cmdGen.Flags().StringP("algo", "a", "eddsa", "Private key algorithm")
	if err := cmdGen.MarkFlagRequired("name"); err != nil {
		log.Fatal().Err(err).Msg("unable to mark name flag as required")
	}
}
