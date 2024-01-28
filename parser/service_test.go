package parser

import (
	"errors"
	"github.com/go-services/source"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gs/fs"
	"os"
	"testing"
)

func newTestFile() {

}

func TestFindServicesWithNoServices(t *testing.T) {
	setup()
	src, _ := source.New(
		`package main
type Test interface {
	Example()
}

type test struct {}
 
func New() Test {
 return &test{}
}
func (c test) Example() {}`,
	)
	file := AnnotatedFile{
		Path: "no_service.go",
		Src:  *src,
	}
	apis, err := FindServices([]AnnotatedFile{file})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(apis), "should be equal")

}
func TestFindApis(t *testing.T) {
	setup()
	src, _ := source.New(`package main
import "context"
//@service()
type Test interface {
	//@http(method="GET", path="/testing")
	Example(ctx context.Context) error
}

type test struct {
}
 
func New() Test {
 return &test{}
}
func (c test) Example(ctx context.Context) {}`,
	)
	file := AnnotatedFile{
		Path: "one_service.go",
		Src:  *src,
	}
	apis, err := FindServices([]AnnotatedFile{file})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(apis), "should be equal")
}

func TestFindServicesWithMultipleServices(t *testing.T) {
	setup()

	src, _ := source.New(`package main
import "context"
//@service()
type Test interface {
	//@http(method="GET", path="/testing")
	Example(ctx context.Context) error
}

type test struct {
}
 
func New() Test {
 return &test{}
}
func (c test) Example(ctx context.Context) {}`,
	)

	src2, _ := source.New(`package main
import "context"
//@service()
type AnotherTest interface {
	//@http(method="GET", path="/testing")
	Example(ctx context.Context) error
}

type test struct {
}
 
func New() AnotherTest {
 return &test{}
}
func (c test) Example(ctx context.Context) {}`,
	)
	file1 := AnnotatedFile{
		Path: "first_service.go",
		Src:  *src,
	}
	file2 := AnnotatedFile{
		Path: "first_service.go",
		Src:  *src2,
	}
	apis, err := FindServices([]AnnotatedFile{file1, file2})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(apis), "should be equal")
}

func TestFindServicesWithDuplicateServices(t *testing.T) {
	setup()

	src, _ := source.New(`package main
    import "context"
    //@service()
    type Test interface {
        //@http(method="GET", path="/testing")
        Example(ctx context.Context) error
    }

    type test struct {
    }

    func New() Test {
        return &test{}
    }
    func (c test) Example(ctx context.Context) {}`,
	)

	file1 := AnnotatedFile{
		Path: "duplicate_service.go",
		Src:  *src,
	}
	file2 := AnnotatedFile{
		Path: "duplicate_service.go",
		Src:  *src,
	}
	_, err := FindServices([]AnnotatedFile{file1, file2})
	assert.Equal(t, errors.New("service `test` is defined more than once"), err)
}

func TestFindServicesWithSameBaseRouteShouldReturnError(t *testing.T) {
	setup()

	src, _ := source.New(`package main
    import "context"
    //@service(base="/testing")
    type Test interface {
        //@http(method="GET", path="/testing")
        Example(ctx context.Context) error
    }

    type test struct {
    }

    func New() Test {
        return &test{}
    }
    func (c test) Example(ctx context.Context) {}`,
	)

	src2, _ := source.New(`package main
    import "context"
    //@service(base="/testing")
    type AnotherTest interface {
        //@http(method="GET", path="/testing")
        Example(ctx context.Context) error
    }

    type test struct {
    }

    func New() AnotherTest {
        return &test{}
    }
    func (c test) Example(ctx context.Context) {}`,
	)
	file1 := AnnotatedFile{
		Path: "first_service.go",
		Src:  *src,
	}
	file2 := AnnotatedFile{
		Path: "second_service.go",
		Src:  *src2,
	}
	_, err := FindServices([]AnnotatedFile{file1, file2})
	assert.Equal(t, errors.New("service `another_test` has the same base route as `test`"), err)
}

func TestParseServiceWithNoConstructorFunction(t *testing.T) {
	setup()

	src, _ := source.New(`package main
    import "context"
    //@service()
    type Test interface {
        //@http(method="GET", path="/testing")
        Example(ctx context.Context) error
    }

    type test struct {
    }

    func (c test) Example(ctx context.Context) {}`,
	)

	file := AnnotatedFile{
		Path: "invalid_constructor.go",
		Src:  *src,
	}

	_, err := parseService(file)
	assert.Equal(t, errors.New("service `test` has no constructor function, each service needs to have a function that returns a new instance."), err)

}

func TestParseServiceWithSameBasePathAndMethod(t *testing.T) {
	setup()
	src, _ := source.New(`package main
        import "context"
        //@service()
        type MyService interface {
            //@http(method="GET", path="/testing")
            FirstMethod(ctx context.Context) error
            //@http(method="GET", path="/testing")
            SecondMethod(ctx context.Context) error
        }

        type test struct {
        }

        func New() MyService {
            return &test{}
        }
        func (c test) FirstMethod(ctx context.Context) {}
        func (c test) SecondMethod(ctx context.Context) {}`,
	)

	file := AnnotatedFile{
		Path: "same_base_path_and_method.go",
		Src:  *src,
	}

	_, err := parseService(file)
	assert.Equal(t, errors.New("endpoint `MyService.SecondMethod` has the same basePath and method as `FirstMethod`"), err)
}

func TestParseEndpointWithInvalidParams(t *testing.T) {
	setup()
	src, _ := source.New(`package main
        import "context"
        //@service()
        type MyService interface {
            //@http(method="GET", path="/testing")
            InvalidParamsMethod(ctx context.Context, extraParam int) error
        }

        type test struct {
        }

        func New() MyService {
            return &test{}
        }
        func (c test) InvalidParamsMethod(ctx context.Context, extraParam int) {}`,
	)

	file := AnnotatedFile{
		Path: "invalid_params.go",
		Src:  *src,
	}

	_, err := parseService(file)
	assert.Equal(t, errors.New("error while parsing endpoint `MyService.InvalidParamsMethod` : request needs to be an exported structure"), err)
}

func TestParseEndpointWithInvalidResults(t *testing.T) {
	setup()
	src, _ := source.New(`package main
        import "context"
        //@service()
        type MyService interface {
            //@http(method="GET", path="/testing")
            InvalidResultsMethod(ctx context.Context) (int, error)
        }

        type test struct {
        }

        func New() MyService {
            return &test{}
        }
        func (c test) InvalidResultsMethod(ctx context.Context) (int, error) { return 0, nil }`,
	)

	file := AnnotatedFile{
		Path: "invalid_results.go",
		Src:  *src,
	}

	_, err := parseService(file)
	assert.Equal(t, errors.New("error while parsing endpoint `MyService.InvalidResultsMethod` : method must return either the error or the response pointer and the error"), err)
}

func TestParseServiceWithValidRequestAndResponse(t *testing.T) {
	_ = afero.NewOsFs().MkdirAll("./.tmp", 0755)
	fs.SetTestFs(afero.NewOsFs())
	defer func() {
		_ = os.Chdir("..")
		_ = fs.DeleteFolder("./.tmp")
	}()
	_ = os.Chdir("./.tmp")
	_ = fs.WriteFile("go.mod", "module test")
	_ = fs.WriteFile("valid_request_response.go", `package test
        import "context"
        //@service()
        type MyService interface {
            //@http(method="GET", path="/testing")
            ValidMethod(ctx context.Context, req Request) (*Response, error)
        }

        type Request struct {
            Field1 string
            Field2 int
        }

        type Response struct {
            Field1 string
            Field2 int
        }

        type test struct {
        }

        func New() MyService {
            return &test{}
        }
        func (c test) ValidMethod(ctx context.Context, req Request) (Response, error) {
            return Response{}, nil
        }`)

	file, _ := ParseFiles(".")

	_, err := parseService(file[0])
	assert.Nil(t, err)
}
