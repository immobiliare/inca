package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.rete.farm/sistemi/inca/util"
)

var cmdGen = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := cmd.Flags().GetString("out")
		if err != nil {
			logrus.WithError(err).Fatalln()
		}
		logrus.WithFields(logrus.Fields{
			"output": output,
		}).Println("gen called")
	},
}

func init() {
	cmdRoot.AddCommand(cmdGen)
	cmdGen.Flags().StringP("out", "o", util.ErrWrap("./")(os.Getwd()), "Output path")
}
