package generate

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"gs/assets"
	"gs/config"
	"gs/fs"
	"gs/parser"
	"os/exec"
	"path"
)

type ServiceGenerator interface {
	Generate() error
}

type EndpointMethodData struct {
	Module           string
	Service          string
	ServiceImport    string
	ServiceInterface string
	Endpoint         parser.Endpoint
}

type serviceGenerator struct {
	services []parser.Service
}

func NewServiceGenerator(services []parser.Service) ServiceGenerator {
	return &serviceGenerator{
		services: services,
	}
}

func generateEndpoints(sv parser.Service, endpointPath string) error {
	definitionsFolder := path.Join(endpointPath, "definitions")
	if exists, _ := fs.Exists(definitionsFolder); !exists {
		_ = fs.CreateFolder(definitionsFolder)
	}
	// add method definition
	for _, ep := range sv.Endpoints {
		err := assets.ParseAndWriteTemplate(
			"services/endpoints/definitions/method.go.tmpl",
			path.Join(
				definitionsFolder,
				fmt.Sprintf("%s.go", ep.PackageName()),
			),
			ep,
		)
		if err != nil {
			return err
		}
		err = assets.ParseAndWriteTemplate(
			"services/endpoints/method.go.tmpl",
			path.Join(
				endpointPath,
				fmt.Sprintf("%s.go", ep.PackageName()),
			),
			EndpointMethodData{
				Module:           sv.Config.Module,
				Service:          sv.FormattedName,
				ServiceImport:    sv.Import,
				ServiceInterface: sv.Interface,
				Endpoint:         ep,
			},
		)
		if err != nil {
			return err
		}
	}
	err := assets.ParseAndWriteTemplate(
		"services/endpoints/options.go.tmpl",
		path.Join(
			endpointPath,
			"options$.go",
		),
		sv,
	)
	if err != nil {
		return err
	}
	err = assets.ParseAndWriteTemplate(
		"services/endpoints/endpoint.go.tmpl",
		path.Join(
			endpointPath,
			"endpoint$.go",
		),
		sv,
	)
	if err != nil {
		return err
	}
	return nil
}

func generateHttpTransport(svc parser.Service, httpTransportPath string) error {
	globalTransportPth := path.Join(genPath(), "transport", "http")
	if exists, _ := fs.Exists(globalTransportPth); !exists {
		_ = fs.CreateFolder(globalTransportPth)
	}
	err := assets.ParseAndWriteTemplate(
		"transport/http/http.go.tmpl",
		path.Join(
			globalTransportPth,
			"http.go",
		),
		nil,
	)
	if err != nil {
		return err
	}

	for _, ep := range svc.Endpoints {
		err = assets.ParseAndWriteTemplate(
			"services/transport/http/method.go.tmpl",
			path.Join(
				httpTransportPath,
				fmt.Sprintf("%s.go", ep.PackageName()),
			),
			EndpointMethodData{
				Module:           svc.Config.Module,
				Service:          svc.FormattedName,
				ServiceImport:    svc.Import,
				ServiceInterface: svc.Interface,
				Endpoint:         ep,
			},
		)
		if err != nil {
			return err
		}
	}
	err = assets.ParseAndWriteTemplate(
		"services/transport/http/http.go.tmpl",
		path.Join(
			httpTransportPath,
			"http$.go",
		),
		svc,
	)
	if err != nil {
		return err
	}

	return nil
}

func generateService(svc parser.Service) error {
	log.Debug("Starting service generation for ", svc.Name)
	pth := path.Join(genPath(), "services", svc.FormattedName)
	epFolder := path.Join(pth, "endpoint")
	httpTransportPath := path.Join(pth, "transport", "http")
	if exists, _ := fs.Exists(epFolder); !exists {
		_ = fs.CreateFolder(epFolder)
	}
	err := generateEndpoints(svc, epFolder)
	if err != nil {
		log.Error("Error generating endpoints for ", svc.Name, ": ", err)
		return err
	}

	err = generateHttpTransport(svc, httpTransportPath)
	if err != nil {
		log.Error("Error generating HTTP transport for ", svc.Name, ": ", err)
		return err
	}

	err = assets.ParseAndWriteTemplate(
		"services/service.go.tmpl",
		path.Join(
			pth,
			"service.go",
		),
		svc,
	)
	if err != nil {
		log.Error("Error generating service template for ", svc.Name, ": ", err)
	}
	log.Debug("Finished service generation for ", svc.Name)
	return err
}

func (g serviceGenerator) Generate() error {
	for _, svc := range g.services {
		err := generateService(svc)
		if err != nil {
			return err
		}
		handlerPath := path.Join(cmdPath(), "app", svc.FormattedName)
		if exists, _ := fs.Exists(handlerPath); !exists {
			_ = fs.CreateFolder(handlerPath)
		}

		handlerPath = path.Join(handlerPath, fmt.Sprintf("%s.go", svc.FormattedName))
		if exists, _ := fs.Exists(handlerPath); !exists {
			err = assets.ParseAndWriteTemplate(
				"cmd/service.go.tmpl",
				handlerPath,
				svc,
			)
			if err != nil {
				return err
			}
		}

		mainHandlerPath := path.Join(cmdPath(), fmt.Sprintf("%s.go", svc.FormattedName))
		if exists, _ := fs.Exists(mainHandlerPath); !exists {
			err = assets.ParseAndWriteTemplate(
				"cmd/cmd.go.tmpl",
				mainHandlerPath,
				svc,
			)
			if err != nil {
				return err
			}
		}
	}
	handlerPath := genPath()
	if exists, _ := fs.Exists(handlerPath); !exists {
		_ = fs.CreateFolder(handlerPath)
	}

	handlerPath = path.Join(handlerPath, "cmd", fmt.Sprintf("%s.go", strcase.ToSnake(config.Get().Module)))
	if exists, _ := fs.Exists(handlerPath); !exists {
		err := assets.ParseAndWriteTemplate(
			"cmd/all.go.tmpl",
			handlerPath,
			g.services,
		)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command("go", "mod", "tidy")
	return cmd.Run()
}
