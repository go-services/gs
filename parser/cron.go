package parser

import (
	"errors"
	"fmt"
	"github.com/go-services/annotation"
	"github.com/go-services/source"
	"github.com/iancoleman/strcase"
	"github.com/robfig/cron/v3"
	"gs/config"
)

type Cron struct {
	Name string

	Config config.GSConfig

	Interface string

	Import string

	Package       string
	FormattedName string

	NewMethod string

	Schedule string
}

func parseCronJob(file AnnotatedFile) (*Cron, error) {
	cnf := config.Get()
	inf, ann := findCronJobInterface(file.Src)
	if inf == nil {
		return nil, nil
	}
	name := strcase.ToSnake(ann.Get("name").String())
	if name == "" {
		name = strcase.ToSnake(inf.Name())
	}

	schedule := ann.Get("schedule").String()
	if schedule == "" {
		return nil, errors.New(fmt.Sprintf("cron %s has no schedule annotation", name))
	}
	if _, err := cron.ParseStandard(schedule); err != nil {
		return nil, errors.New(fmt.Sprintf("cron %s has an invalid schedule annotation", name))
	}

	var newMth *source.Function
	for _, fn := range file.Src.Functions() {
		if isExported(fn.Name()) && len(fn.Results()) == 1 && fn.Results()[0].Type.Qualifier == inf.Name() {
			newMth = &fn
			break
		}
	}
	if newMth == nil {
		return nil, errors.New(fmt.Sprintf("cron %s has no constructor function, each cron needs to have a function that returns a new instance", name))
	}
	var runMth *source.InterfaceMethod
	for _, method := range inf.Methods() {
		if method.Name() == "Run" && len(method.Params()) == 0 && len(method.Results()) == 0 {
			runMth = &method
			break
		}
	}
	if runMth == nil {
		return nil, errors.New(fmt.Sprintf("cron %s has no Run method, each cron needs to have a method called Run", name))
	}
	cron := &Cron{
		Name:          name,
		Config:        *cnf,
		Interface:     inf.Name(),
		NewMethod:     newMth.Name(),
		Import:        getPackageImport(cnf.Module, file.Path, file.Src.Package()),
		Package:       file.Src.Package(),
		FormattedName: strcase.ToSnake(name),
		Schedule:      schedule,
	}
	return cron, nil
}

func findCronJobInterface(src source.Source) (*source.Interface, *annotation.Annotation) {
	for _, inf := range src.Interfaces() {
		annotations := findAnnotations("cron", inf.Annotations())
		if len(annotations) > 0 {
			if len(annotations) > 1 {
				log.Warnf("Interface `%s` has more than one cron annotation last one will be used", inf.Name())
			}
			return &inf, &annotations[len(annotations)-1]
		}
	}
	return nil, nil
}
func FindCronJobs(files []AnnotatedFile) (jobs []Cron, err error) {
	for _, file := range files {
		job, err := parseCronJob(file)
		if err != nil {
			return nil, err
		}
		if job == nil {
			continue
		}
		for _, j := range jobs {
			if j.Name == job.Name {
				return nil, errors.New(fmt.Sprintf("cron %s already exists", job.Name))
			}
		}
		jobs = append(jobs, *job)
	}
	return jobs, nil
}
