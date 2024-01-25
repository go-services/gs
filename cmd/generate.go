package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gs/generate"
)

var generateCmd = &cobra.Command{
	Use: "generate",
	Aliases: []string{
		"gen",
		"g",
	},
	Short: "Generate services",
	RunE: func(cmd *cobra.Command, args []string) error {
		if b, _ := cmd.Flags().GetBool("debug"); b {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return generate.Generate()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
