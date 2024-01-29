package cmd

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gs/config"
	"gs/generate"
	"os/exec"
)

var pluginsCmd = &cobra.Command{
	Use: "plugins",
	Aliases: []string{
		"p",
	},
	Short: "Plugins commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var addPluginCmd = &cobra.Command{
	Use: "add",
	Aliases: []string{
		"a",
	},
	Short: "Add plugin",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid number of arguments")
		}
		v := args[0]
		if v != "sst" {
			return errors.New("invalid plugin")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !config.Exists() {
			return errors.New("gs.yaml does not exist")
		}
		cnf := config.Get()
		if cnf.SST != nil {
			logrus.Info("SST plugin already exists")
			return nil
		}

		cnf.SST = &config.SSTConfig{
			StacksPath: "stacks",
		}
		err := cnf.Write()
		if err != nil {
			return err
		}
		cnf.Reload()
		logrus.Info("SST plugin added")
		err = generate.Generate()
		if err != nil {
			return err
		}
		npmCmd := exec.Command("npm", "install")
		npmCmd.Stderr = cmd.ErrOrStderr()
		npmCmd.Stdout = cmd.OutOrStdout()
		return npmCmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(pluginsCmd)
	pluginsCmd.AddCommand(addPluginCmd)
}
