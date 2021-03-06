package cmd

import (
	"gs/service"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use: "service",
	Aliases: []string{
		"svc",
		"s",
	},
	Short: "Create a new service in the project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if b, _ := cmd.Flags().GetBool("debug"); b {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return service.GenerateNew(args[0])
	},
}

func init() {
	newCmd.AddCommand(serviceCmd)
}
