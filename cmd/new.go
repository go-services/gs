package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gs/assets"
	"gs/fs"
	"gs/generate"
	"os"
	"path"
	"runtime"
	"strings"
)

var newCmd = &cobra.Command{
	Use: "new",
	Aliases: []string{
		"n",
	},
	Short: "New",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var projectCmd = &cobra.Command{
	Use: "project",
	Aliases: []string{
		"p",
	},
	Args:  cobra.ExactArgs(1),
	Short: "New project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if b, _ := cmd.Flags().GetBool("debug"); b {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return generateProject(args[0])
	},
}

func init() {
	newCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(newCmd)
}

func generateProject(name string) error {
	formattedName := strcase.ToKebab(name)
	if exists, _ := fs.Exists(formattedName); exists {
		logrus.Errorf("Folder %s already exists", formattedName)
		return nil
	}
	err := assets.ParseAndWriteTemplate(
		"project/package.json.tmpl",
		path.Join(formattedName, "package.json"),
		map[string]string{
			"Name": formattedName,
		},
	)
	if err != nil {
		return err
	}

	_ = fs.CreateFolder(path.Join(formattedName, "example"))
	err = assets.ParseAndWriteTemplate(
		"project/example/service.go.tmpl",
		path.Join(formattedName, "example", "service.go"),
		map[string]string{
			"Module": strcase.ToSnake(formattedName),
		},
	)
	if err != nil {
		return err
	}

	err = assets.ParseAndWriteTemplate(
		"project/go.mod.tmpl",
		path.Join(formattedName, "go.mod"),
		map[string]string{
			"Module":  strcase.ToSnake(formattedName),
			"Version": strings.TrimPrefix(runtime.Version(), "go"),
		},
	)
	if err != nil {
		return err
	}

	err = os.Chdir(formattedName)
	if err != nil {
		return err
	}

	log.Infof("Generating project %s", formattedName)

	err = generate.Generate()

	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Project %s created", formattedName)
	log.Infof("Run `cd %s && go run gen/cmd/%s.go` to start the app or gs watch to run and watch for changes", formattedName, strcase.ToSnake(formattedName))

	return nil
}
