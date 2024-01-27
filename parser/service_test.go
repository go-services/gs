package parser

import (
	"errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gs/fs"
	"os"
	"testing"
)

func TestFindServicesWithNoServices(t *testing.T) {
	fs.SetTestFs(afero.NewOsFs())
	_ = os.Chdir("../_testdata/empty")
	files, _ := ParseFiles(".")
	apis, err := FindServices(files)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(apis), "should be equal")
	_ = os.Chdir("../../parser")

}
func TestFindApis(t *testing.T) {
	fs.SetTestFs(afero.NewOsFs())
	_ = os.Chdir("../_testdata/one")
	files, _ := ParseFiles(".")
	apis, _ := FindServices(files)
	assert.Equal(t, 1, len(apis), "should be equal")
	_ = os.Chdir("../../parser")
}

func TestFindServicesWithMultipleServices(t *testing.T) {
	fs.SetTestFs(afero.NewOsFs())
	_ = os.Chdir("../_testdata/multiple")
	files, _ := ParseFiles(".")
	apis, err := FindServices(files)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(apis), "should be equal")
	_ = os.Chdir("../../parser")
}

func TestFindServicesWithDuplicateServices(t *testing.T) {
	fs.SetTestFs(afero.NewOsFs())
	_ = os.Chdir("../_testdata/duplicate")
	files, _ := ParseFiles(".")
	_, err := FindServices(files)
	assert.Equal(t, errors.New("service `duplicate` is defined more than once"), err)
	_ = os.Chdir("../../parser")
}

func TestFindServicesWithSameBaseRoute(t *testing.T) {
	fs.SetTestFs(afero.NewOsFs())
	_ = os.Chdir("../_testdata/same_base_route")
	files, _ := ParseFiles(".")
	_, err := FindServices(files)
	assert.Equal(t, errors.New("service `service_2` has the same base route as `service_1`"), err)
	_ = os.Chdir("../../parser")
}

func TestParseServiceWithNoConstructorFunction(t *testing.T) {
	fs.SetTestFs(afero.NewOsFs())
	_ = os.Chdir("../_testdata/no_constructor")
	files, _ := ParseFiles(".")
	file := files[0]
	_, err := parseService(file)
	assert.Equal(t, errors.New("service `no_constructor` has no constructor function, each service needs to have a function that returns a new instance."), err)
	_ = os.Chdir("../../parser")
}

func TestParseServiceWithSameBasePathAndMethod(t *testing.T) {
	fs.SetTestFs(afero.NewOsFs())
	_ = os.Chdir("../_testdata/same_base_path_and_method")
	files, _ := ParseFiles(".")
	file := files[0]
	_, err := parseService(file)
	assert.Equal(t, errors.New("endpoint `XyzService.MyMethod` has the same basePath and method as `SecondMethod`"), err)
	_ = os.Chdir("../../parser")
}

func TestParseEndpointWithInvalidParams(t *testing.T) {
	fs.SetTestFs(afero.NewOsFs())
	_ = os.Chdir("../_testdata/invalid_params")
	files, _ := ParseFiles(".")
	file := files[0]
	_, err := parseService(file)
	assert.Equal(t, errors.New("error while parsing endpoint `XyzService.MyMethod` : method must except either the context or the context and the request struct"), err)
	_ = os.Chdir("../../parser")
}

func TestParseEndpointWithInvalidResults(t *testing.T) {
	fs.SetTestFs(afero.NewOsFs())
	_ = os.Chdir("../_testdata/invalid_results")
	files, _ := ParseFiles(".")
	file := files[0]
	_, err := parseService(file)
	assert.Equal(t, errors.New("error while parsing endpoint `XyzService.MyMethod` : method must return either the error or the response pointer and the error"), err)
	_ = os.Chdir("../../parser")
}
