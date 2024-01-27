package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-services/code"

	"github.com/go-services/annotation"
)

type paramType string
type encodeFormat string

func (v encodeFormat) String() string {
	return string(v)
}

const (
	URL   paramType = "URL"
	QUERY paramType = "QUERY"
	BODY  paramType = "BODY"
)

const (
	JSON encodeFormat = "JSON"
	XML  encodeFormat = "XML"
	FORM encodeFormat = "FORM"
)

type ParamParser struct {
	Fn      string
	NoError bool
}

var typeFuncMap = map[string]*ParamParser{
	"[]string": {
		Fn:      "StringToStringArray",
		NoError: true,
	},
	"int": {
		Fn: "StringToInt",
	},
	"int64": {
		Fn: "StringToInt64",
	},
	"[]int": {
		Fn: "StringToIntArray",
	},
	"[]int64": {
		Fn: "StringToInt64Array",
	},
	"float64": {
		Fn: "StringToFloat64",
	},
	"[]float64": {
		Fn: "StringToFloat64Array",
	},
	"float32": {
		Fn: "StringToFloat32",
	},
	"[]float32": {
		Fn: "StringToFloat32Array",
	},
	"bool": {
		Fn: "StringToBool",
	},
}

type HttpRequestParam struct {
	// this is the field Name
	Field string
	// this is the Name given in the url param or query param
	Name string
	// this is the field type
	Type code.Type
	// is this parameter optional
	Required bool
	// this tells us of it is a URL param or a Query param
	ParamType paramType
	// parameter parse function
	Parser *ParamParser
}

type HttpRequest struct {
	// the format the data is
	Format encodeFormat
	// if the request has any query params
	HasUrl bool
	// if the request has body portion
	HasBody bool
	// all the extra params
	Params []HttpRequestParam
}

type HttpMethodRoute struct {
	Name    string
	Methods []string
	Route   string
}

type HttpTransport struct {
	Request        *HttpRequest
	ResponseFormat encodeFormat
	MethodRoutes   []HttpMethodRoute
}

func parseHttpTransport(serviceRoute string, endpoint Endpoint) (*HttpTransport, error) {
	httpAnnotations := findAnnotations("http", endpoint.Annotations)
	if len(httpAnnotations) == 0 {
		return nil, nil
	}

	methodRoutes, err := parseMethodRoutes(serviceRoute, httpAnnotations[0])

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error while parsing endpoint `%s` : %s", endpoint.Name, err.Error()))
	}
	return &HttpTransport{
		MethodRoutes:   methodRoutes,
		Request:        parseHttpRequest(endpoint),
		ResponseFormat: encodeFormat(httpResponseFormat(httpAnnotations[0].Get("response").String())),
	}, nil
}

func httpResponseFormat(format string) encodeFormat {
	if format == "" {
		return JSON
	}
	switch encodeFormat(strings.ToUpper(format)) {
	case XML:
		return XML
	case FORM:
		return FORM
	default:
		return JSON
	}
}
func parseHttpRequest(endpoint Endpoint) *HttpRequest {
	if endpoint.Request == nil {
		return nil
	}

	httpAnnotations := findAnnotations("http", endpoint.Annotations)
	httpEncode := httpAnnotations[0]
	annotationFormat := httpEncode.Get("request").String()
	format := JSON
	if annotationFormat != "" {
		switch encodeFormat(strings.ToUpper(annotationFormat)) {
		case XML:
			format = XML
		case FORM:
			format = FORM
		case JSON:
			format = JSON
		default:
			log.WithField("endpoint", endpoint.Name).Info("The request format is not supported `json` will be used as default")
		}
	}
	request := &HttpRequest{
		Format: format,
	}
	parseHttpRequestParams(endpoint.Request, request)
	return request
}

func parseHttpRequestParams(req *code.Struct, request *HttpRequest) {
	for _, field := range req.Fields {
		if !isExported(field.Name) || field.Tags == nil {
			continue
		}

		gsUrl := getTag("url", *field.Tags)
		gsQuery := getTag("query", *field.Tags)
		gsBody := getTag("body", *field.Tags)

		tp := field.Type.String()

		if gsUrl != "" {
			if !isUrlTypeSupported(tp) {
				log.WithField("field", field.Name).WithField("type", field.Type.String()).Warn("Field type not supported for url")
				continue
			}
			var parser *ParamParser = nil
			if tp != "string" {
				parser = typeFuncMap[tp]
			}
			name, required := getParameter(gsUrl)
			request.Params = append(request.Params, HttpRequestParam{
				Field:     field.Name,
				Name:      name,
				Type:      field.Type,
				Required:  required,
				ParamType: URL,
				Parser:    parser,
			})
			request.HasUrl = true
		}
		if gsQuery != "" {
			if !isQueryTypeSupported(tp) {
				log.WithField("field", field.Name).WithField("type", field.Type.String()).Warn("Field type not supported for query")
				continue
			}
			var parser *ParamParser = nil
			if tp != "string" {
				parser = typeFuncMap[tp]
			}

			name, required := getParameter(gsQuery)
			request.Params = append(request.Params, HttpRequestParam{
				Field:     field.Name,
				Name:      name,
				Type:      field.Type,
				Required:  required,
				ParamType: QUERY,
				Parser:    parser,
			})
		}
		if gsBody != "" {
			name, required := getParameter(gsBody)
			format := JSON
			switch encodeFormat(strings.ToUpper(name)) {
			case XML:
				format = XML
			case FORM:
				format = FORM
			case JSON:
				format = JSON
			default:
				log.WithField("endpoint", field.Name).Info("The request format is not supported `json` will be used as default")
			}
			request.Params = append(request.Params, HttpRequestParam{
				Field:     field.Name,
				Name:      string(format),
				Required:  required,
				ParamType: BODY,
			})
			request.HasBody = true
		}
	}
	return
}

func parseMethodRoutes(serviceRoute string, httpAnnotation annotation.Annotation) (routes []HttpMethodRoute, err error) {
	keepTrailingSlash := httpAnnotation.Get("keepTrailingSlash").Bool()
	var methodsPrepared []string
	// for now we won't support multiple methods for the same path, if I find that it is used it can be enabled.
	//for _, method := range strings.Split(httpAnnotation.Get("method").String(), ",") {
	//	methodsPrepared = append(methodsPrepared, strings.ToUpper(strings.TrimSpace(method)))
	//}
	method := strings.ToUpper(strings.TrimSpace(httpAnnotation.Get("method").String()))
	// check if the method is valid
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" && method != "PATCH" && method != "OPTIONS" && method != "HEAD" {
		return nil, errors.New(fmt.Sprintf("method `%s` is not valid", method))
	}
	methodsPrepared = append(
		methodsPrepared,
		method,
	)

	path := httpAnnotation.Get("path").String()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	methodRoute := HttpMethodRoute{
		Name:    httpAnnotation.Get("Name").String(),
		Methods: methodsPrepared,
		Route:   serviceRoute + path,
	}
	routes = append(
		routes,
		methodRoute,
	)
	if !keepTrailingSlash {
		if strings.HasSuffix(methodRoute.Route, "/") {
			path = strings.TrimSuffix(path, "/")
		} else {
			path += "/"
		}
		methodRoute.Route = serviceRoute + path
		routes = append(
			routes,
			methodRoute,
		)
	}
	return
}

func isQueryTypeSupported(tp string) bool {
	var supportedQueryTypes = []string{
		"string",
		"[]string",
		"int",
		"int64",
		"[]int",
		"[]int64",
		"bool",
		"float32",
		"[]float32",
		"float64",
		"[]float64",
	}
	for _, supportedType := range supportedQueryTypes {
		if supportedType == tp {
			return true
		}
	}
	return false
}

func isUrlTypeSupported(tp string) bool {
	var supportedUrlTypes = []string{"string", "int", "int64", "float32", "float64"}
	for _, supportedType := range supportedUrlTypes {
		if supportedType == tp {
			return true
		}
	}
	return false
}
