package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	cmdRoot = &cobra.Command{
		Use:   "inca",
		Short: "Inca is an INternal CA manager for local CAs as well as external ones",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// logrus.SetFormatter(&logrus.JSONFormatter{})
			// if debug, err := cmd.Flags().GetBool("debug"); err == nil && debug {
			// 	logrus.SetLevel(logrus.DebugLevel)
			// }
		},
	}
)

func Execute() {
	if cmdRoot.Execute() != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.inca.yaml)")
	// cmdRoot.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmdRoot.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode")
}
