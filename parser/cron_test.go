package parser

import (
	"errors"
	"github.com/go-services/source"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCronJob_WithValidFile(t *testing.T) {
	setup()
	src, _ := source.New("package main\n// @cron(name='valid', schedule='* * * * *')\ntype MyCron interface {\nRun()\n}\n type myCron struct {\n}\n \nfunc New() MyCron {\n return &myCron{}\n}\nfunc (c myCron) Run() {}")
	file := AnnotatedFile{
		Path: "valid.go",
		Src:  *src,
	}
	cron, err := parseCronJob(file)

	assert.Nil(t, err, "should be nil")
	assert.Equal(t, "valid", cron.Name, "should be equal")
	assert.Equal(t, "* * * * *", cron.Schedule, "should be equal")
}

func TestParseCronJob_WithMissingSchedule(t *testing.T) {
	setup()
	src, _ := source.New("package main\n// @cron(name='missing_schedule')\ntype MyCron interface {\nRun()\n}\n type myCron struct {\n}\n \nfunc New() MyCron {\n return &myCron{}\n}\nfunc (c myCron) Run() {}")
	file := AnnotatedFile{
		Path: "missing_schedule.go",
		Src:  *src,
	}
	_, err := parseCronJob(file)

	assert.NotNil(t, err, "should not be nil")
	assert.Equal(t, errors.New("cron missing_schedule has no schedule annotation"), err, "should be equal")
}

func TestParseCronJob_WithInvalidSchedule(t *testing.T) {
	setup()
	src, _ := source.New("package main\n// @cron(name='invalid_schedule', schedule='invalid')\ntype MyCron interface {\nRun()\n}\n type myCron struct {\n}\n \nfunc New() MyCron {\n return &myCron{}\n}\nfunc (c myCron) Run() {}")
	file := AnnotatedFile{
		Path: "invalid_schedule.go",
		Src:  *src,
	}
	_, err := parseCronJob(file)

	assert.NotNil(t, err, "should not be nil")
	assert.Equal(t, errors.New("cron invalid_schedule has an invalid schedule annotation"), err, "should be equal")
}

func TestParseCronJob_WithMissingConstructor(t *testing.T) {
	setup()
	src, _ := source.New("package main\n// @cron(name='missing_constructor', schedule='* * * * *')\ntype MyCron interface {\nRun()\n}\n type myCron struct {\n}\nfunc (c myCron) Run() {}")
	file := AnnotatedFile{
		Path: "missing_constructor.go",
		Src:  *src,
	}
	_, err := parseCronJob(file)

	assert.NotNil(t, err, "should not be nil")
	assert.Equal(t, errors.New("cron missing_constructor has no constructor function, each cron needs to have a function that returns a new instance"), err, "should be equal")
}

func TestParseCronJob_WithMissingRunMethod(t *testing.T) {
	setup()
	src, _ := source.New("package main\n// @cron(name='missing_run', schedule='* * * * *')\ntype MyCron interface {\n}\n type myCron struct {\n}\nfunc New() MyCron {\n return &myCron{}\n}\nfunc (c myCron) Run() {}")
	file := AnnotatedFile{
		Path: "missing_run.go",
		Src:  *src,
	}
	_, err := parseCronJob(file)

	assert.NotNil(t, err, "should not be nil")
	assert.Equal(t, errors.New("cron missing_run has no Run method, each cron needs to have a method called Run"), err, "should be equal")
}

func TestFindCronJobs_WithValidFiles(t *testing.T) {
	setup()
	src, _ := source.New("package main\n// @cron(name='valid1', schedule='* * * * *')\ntype MyCron interface {\nRun()\n}\n type myCron struct {\n}\n \nfunc New() MyCron {\n return &myCron{}\n}\nfunc (c myCron) Run() {}")
	src2, _ := source.New("package main\n// @cron(name='valid2', schedule='* * * * *')\ntype MyCron interface {\nRun()\n}\n type myCron struct {\n}\n \nfunc New() MyCron {\n return &myCron{}\n}\nfunc (c myCron) Run() {}")
	files := []AnnotatedFile{
		{
			Path: "valid1.go",
			Src:  *src,
		},
		{
			Path: "valid2.go",
			Src:  *src2,
		},
	}
	jobs, err := FindCronJobs(files)

	assert.Nil(t, err, "should be nil")
	assert.Len(t, jobs, 2, "should be equal")
}

func TestFindCronJobs_WithDuplicateCronJobs(t *testing.T) {
	setup()
	src, _ := source.New("package main\n// @cron(name='duplicate', schedule='* * * * *')\ntype MyCron interface {\nRun()\n}\n type myCron struct {\n}\n \nfunc New() MyCron {\n return &myCron{}\n}\nfunc (c myCron) Run() {}")
	src2, _ := source.New("package main\n// @cron(name='duplicate', schedule='* * * * *')\ntype MyCron interface {\nRun()\n}\n type myCron struct {\n}\n \nfunc New() MyCron {\n return &myCron{}\n}\nfunc (c myCron) Run() {}")
	files := []AnnotatedFile{
		{
			Path: "duplicate1.go",
			Src:  *src,
		},
		{
			Path: "duplicate2.go",
			Src:  *src2,
		},
	}
	_, err := FindCronJobs(files)

	assert.NotNil(t, err, "should not be nil")
	assert.Equal(t, errors.New("cron duplicate already exists"), err, "should be equal")
}
