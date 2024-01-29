package generate

import (
	"github.com/sirupsen/logrus"
	"gs/config"
	"gs/fs"
	"gs/parser"
	"os/exec"
)

func Generate() error {
	log.Info("Generating services")
	cnf := config.Get()
	if cnf.Module == "" {
		logrus.Error("Not in the root of the module")
		return nil
	}

	_ = fs.DeleteFolder(cnf.Paths.Gen)

	files, err := parser.ParseFiles(".")
	if err != nil {
		return err
	}
	services, err := parser.FindServices(files)
	if err != nil {
		return err
	}
	jobs, err := parser.FindCronJobs(files)
	if err != nil {
		return err
	}

	err = CommonFiles()
	if err != nil {
		return err
	}

	svcGen := NewServiceGenerator(services)
	err = svcGen.Generate()
	if err != nil {
		return err
	}

	if cnf.SST != nil {
		sstGen := NewSSTPlugin(services, jobs)
		err = sstGen.Generate()
		if err != nil {
			return err
		}
	}

	err = LocalImplementation(services, jobs)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "tidy")
	return cmd.Run()
}
