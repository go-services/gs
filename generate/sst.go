package generate

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gs/assets"
	"gs/config"
	"gs/fs"
	"gs/parser"
	"gs/utils"
	"path"
)

type sstPlugin struct {
	services []parser.Service
	cronJobs []parser.Cron
}

func NewSSTPlugin(services []parser.Service, cronJobs []parser.Cron) Generator {
	return &sstPlugin{
		services: services,
		cronJobs: cronJobs,
	}
}

func (s *sstPlugin) Generate() error {
	cnf := config.Get()
	if cnf.SST == nil {
		return errors.New("sst config is not set")
	}
	for _, svc := range s.services {
		lambdaHandler := path.Join(genPath(), "cmd", "sst", "services", svc.FormattedName)
		if exists, _ := fs.Exists(lambdaHandler); !exists {
			_ = fs.CreateFolder(lambdaHandler)
		}

		lambdaHandler = path.Join(lambdaHandler, "main.go")
		if exists, _ := fs.Exists(lambdaHandler); !exists {
			err := assets.ParseAndWriteTemplate(
				"sst/lambda/service.go.tmpl",
				lambdaHandler,
				svc,
			)
			if err != nil {
				return err
			}
		}
	}

	stacksPath := cnf.SST.StacksPath
	if stacksPath == "" {
		stacksPath = "stacks"
	}
	genStack := path.Join(stacksPath, "gen")
	if exists, _ := fs.Exists(genStack); !exists {
		if e, _ := fs.Exists(genStack); !e {
			_ = fs.CreateFolder(genStack)
		}
	}

	var validCronJobs []parser.Cron
	for _, cron := range s.cronJobs {
		if utils.IsValidAwsCron(cron.Schedule) {
			validCronJobs = append(validCronJobs, cron)
		} else {
			logrus.Warnf("Cron %s has an invalid schedule for aws cloudwatch, check https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-cron-expressions.html for more info", cron.Name)
		}
	}
	for _, cron := range validCronJobs {
		lambdaHandler := path.Join(genPath(), "cmd", "sst", "cron", cron.FormattedName)
		if exists, _ := fs.Exists(lambdaHandler); !exists {
			_ = fs.CreateFolder(lambdaHandler)
		}

		lambdaHandler = path.Join(lambdaHandler, "main.go")
		if exists, _ := fs.Exists(lambdaHandler); !exists {
			err := assets.ParseAndWriteTemplate(
				"sst/lambda/cron.go.tmpl",
				lambdaHandler,
				cron,
			)
			if err != nil {
				return err
			}
		}
	}
	err := assets.ParseAndWriteTemplate(
		"sst/stacks/gen/gen.ts.tmpl",
		path.Join(genStack, "index.ts"),
		map[string]interface{}{
			"Services": s.services,
			"Jobs":     validCronJobs,
		},
	)
	if err != nil {
		return err
	}

	if exists, _ := fs.Exists("package.json"); !exists {
		err = assets.ParseAndWriteTemplate(
			"sst/package.json.tmpl",
			"package.json",
			nil,
		)
		if err != nil {
			return err
		}
	}

	if exists, _ := fs.Exists("tsconfig.json"); !exists {
		err = assets.ParseAndWriteTemplate(
			"sst/tsconfig.json.tmpl",
			"tsconfig.json",
			nil,
		)
		if err != nil {
			return err
		}
	}

	if exists, _ := fs.Exists("sst.config.ts"); !exists {
		err = assets.ParseAndWriteTemplate(
			"sst/sst.config.ts.tmpl",
			"sst.config.ts",
			map[string]interface{}{
				"Module": cnf.Module,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
