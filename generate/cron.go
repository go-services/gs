package generate

import (
	"gs/assets"
	"gs/fs"
	"gs/parser"
	"path"
)

type cronGenerator struct {
	cronJobs []parser.Cron
}

func NewCronGenerator(cronJobs []parser.Cron) Generator {
	return &cronGenerator{
		cronJobs: cronJobs,
	}
}

func (c cronGenerator) Generate() error {
	for _, cr := range c.cronJobs {
		handlerPath := path.Join(cmdPath(), "jobs")
		if exists, _ := fs.Exists(handlerPath); !exists {
			_ = fs.CreateFolder(handlerPath)
		}

		handlerPath = path.Join(handlerPath, cr.FormattedName, "cmd", "main.go")
		if exists, _ := fs.Exists(handlerPath); !exists {
			err := assets.ParseAndWriteTemplate(
				"cmd/cron/cmd.go.tmpl",
				handlerPath,
				cr,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
