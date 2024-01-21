package cmd

import (
	"gs/watch"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch is used to hot reload your microservices",
	RunE: func(cmd *cobra.Command, args []string) error {
		if b, _ := cmd.Flags().GetBool("debug"); b {
			logrus.SetLevel(logrus.DebugLevel)
		}
		watch.Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
