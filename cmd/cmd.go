package cmd

import (
	"github.com/ckeyer/logrus"
	"github.com/spf13/cobra"
)

var (
	debug bool

	rootCmd = &cobra.Command{
		Use:   "attack",
		Short: "",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			logrus.Debugf("debug ?: %v", debug)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "debug")
}

func Execute() {
	rootCmd.Execute()
}
