package cmd

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/util"
)

var cmdGen = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		var req = pki.NewRequest()

		output, err := cmd.Flags().GetString("out")
		if err != nil {
			logrus.WithError(err).Fatalln()
		}

		names, err := cmd.Flags().GetStringArray("name")
		if err != nil {
			logrus.WithError(err).Fatalln()
		}
		req.Hosts = names

		algo, err := cmd.Flags().GetString("algo")
		if err != nil {
			logrus.WithError(err).Fatalln()
		}
		req.Algo = map[string]int{"eddsa": pki.EDDSA, "rsa": pki.RSA}[algo]
		logrus.WithFields(logrus.Fields{
			"names":    strings.Join(req.Hosts, ","),
			"duration": req.Duration,
			"algo":     algo,
		}).Infoln("generating certificate")
		crt, key, err := pki.New(req)
		if err != nil {
			logrus.WithError(err).Fatalln()
		}

		logrus.WithField("output", output).Println("persisting certificate")
		if err := pki.Export(crt, key, output); err != nil {
			logrus.WithError(err).Fatalln()
		}

		logrus.WithField("output", output).Println("certificate created")
	},
}

func init() {
	cmdRoot.AddCommand(cmdGen)
	cmdGen.Flags().StringP("out", "o", util.ErrWrap("./")(os.Getwd()), "Output path")
	cmdGen.Flags().StringArrayP("name", "n", []string{}, "Certificate names")
	cmdGen.Flags().StringP("algo", "a", "eddsa", "Private key algorithm")
	if err := cmdGen.MarkFlagRequired("name"); err != nil {
		logrus.WithError(err).Fatalln()
	}
}
