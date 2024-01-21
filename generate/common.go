package generate

import (
	"gs/assets"
	"gs/config"
	"gs/fs"
	"path"
	"strings"
)

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
	if cnf.Paths.Cmd != "" {
		cp = cnf.Paths.Cmd
	}
	return cp
}

func Common() error {
	log.Debug("Starting Common function")
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
	log.Debug("Finished Common function successfully")
	return nil
}
