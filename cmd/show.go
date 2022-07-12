package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.rete.farm/sistemi/inca/pki"
)

var showGen = &cobra.Command{
	Use:   "show [certificate] [key]",
	Short: "Pretty print a certificate",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			logrus.Fatalln("wrong arguments")
		}
		var (
			crtPath = args[0]
			keyPath = args[1]
		)

		logrus.WithField("path", keyPath).Println("parsing key")
		key, err := pki.ParseKey(keyPath)
		if err != nil {
			logrus.WithError(err).Fatalln()
		}
		logrus.WithField("public", key.Public()).Println("key parsed")

		logrus.WithField("path", crtPath).Println("parsing certificate")
		crt, err := pki.Parse(crtPath)
		if err != nil {
			logrus.WithError(err).Fatalln()
		}
		logrus.WithField("names", crt.DNSNames).Println("certificate parsed")
	},
}

func init() {
	cmdRoot.AddCommand(showGen)
}
