package generate

import (
	"gs/assets"
	"gs/config"
	"gs/fs"
	"gs/parser"
	"path"
	"strings"
)

type Generator interface {
	Generate() error
}

func genPath() string {
	cnf := config.Get()
	gp := "gen"
	if cnf.Paths.Gen != "" {
		gp = cnf.Paths.Gen
	}
	return gp
}

func cmdPath() string {
	cnf := config.Get()
	cp := "cmd"
	if cnf.Paths.Config != "" {
		cp = cnf.Paths.Config
	}
	return cp
}

func CommonFiles() error {
	log.Debug("Starting CommonFiles function")
	gp := genPath()
	log.Debug("Generated path: ", gp)
	if exists, _ := fs.Exists(gp); !exists {
		log.Debug("Path does not exist, creating folder: ", gp)
		_ = fs.CreateFolder(gp)
	}

	commonFiles := []string{
		"errors/errors.go.tmpl",
		"errors/http.go.tmpl",
		"utils/utils.go.tmpl",
	}

	for _, file := range commonFiles {
		log.Debug("Processing file: ", file)
		err := assets.ParseAndWriteTemplate(file, path.Join(gp, strings.TrimSuffix(file, ".tmpl")), nil)
		if err != nil {
			log.Error("Error processing file: ", file, " Error: ", err)
			return err
		}
	}
	log.Debug("Finished CommonFiles function successfully")
	return nil
}

func LocalImplementation(services []parser.Service, cronJobs []parser.Cron) error {
	handlerPath := genPath()
	if exists, _ := fs.Exists(handlerPath); !exists {
		_ = fs.CreateFolder(handlerPath)
	}

	handlerPath = path.Join(handlerPath, "cmd", "local", "main.go")
	if exists, _ := fs.Exists(handlerPath); !exists {
		err := assets.ParseAndWriteTemplate(
			"cmd/all.go.tmpl",
			handlerPath,
			map[string]interface{}{
				"Services": services,
				"CronJobs": cronJobs,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
