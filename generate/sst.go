package generate

import (
	"errors"
	"gs/assets"
	"gs/config"
	"gs/fs"
	"gs/parser"
	"path"
)

type SSTPlugin interface {
	Generate() error
}

type sstPlugin struct {
	services []parser.Service
}

func NewSSTPlugin(services []parser.Service) SSTPlugin {
	return &sstPlugin{
		services: services,
	}
}

func (s *sstPlugin) Generate() error {
	cnf := config.Get()
	if cnf.SST == nil {
		return errors.New("sst config is not set")
	}
	for _, svc := range s.services {
		lambdaHandler := path.Join(cmdPath(), svc.Name)
		if exists, _ := fs.Exists(lambdaHandler); !exists {
			_ = fs.CreateFolder(lambdaHandler)
		}

		lambdaHandler = path.Join(lambdaHandler, "lambda.go")
		if exists, _ := fs.Exists(lambdaHandler); !exists {
			err := assets.ParseAndWriteTemplate(
				"cmd/lambda/service.go.tmpl",
				lambdaHandler,
				svc,
			)
			if err != nil {
				return err
			}
		}
	}

	stacksPath := cnf.SST.Path
	if stacksPath == "" {
		stacksPath = "stacks"
	}
	genStack := path.Join(stacksPath, "gen")
	if exists, _ := fs.Exists(genStack); !exists {
		if e, _ := fs.Exists(genStack); !e {
			_ = fs.CreateFolder(genStack)
		}
	}
	err := assets.ParseAndWriteTemplate(
		"project/stacks/gen/gen.ts.tmpl",
		path.Join(genStack, "index.ts"),
		s.services,
	)
	if err != nil {
		return err
	}
	return nil
}
