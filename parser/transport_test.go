package parser

import (
	"github.com/go-services/annotation"
	"github.com/go-services/code"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHttpTransportWithValidAnnotations(t *testing.T) {
	ann, _ := annotation.Parse("@http(method=\"GET\", path=\"/test\", request=\"json\", response=\"json\")")
	endpoint := Endpoint{
		Annotations: []annotation.Annotation{
			*ann,
		},
		Request: code.NewStruct("Abc"),
	}
	transport, err := parseHttpTransport("/service", endpoint)
	assert.Nil(t, err)
	assert.NotNil(t, transport)
	assert.Equal(t, JSON, transport.Request.Format)
	assert.Equal(t, JSON, transport.ResponseFormat)
	assert.Equal(t, "/service/test", transport.MethodRoutes[0].Route)
}

func TestParseHttpTransportWithNoAnnotations(t *testing.T) {
	endpoint := Endpoint{}
	transport, err := parseHttpTransport("/service", endpoint)
	assert.Nil(t, err)
	assert.Nil(t, transport)
}

func TestParseHttpTransportWithInvalidAnnotations(t *testing.T) {
	ann, _ := annotation.Parse("@http(method=\"INVALID\", path=\"/test\", request=\"json\", response=\"json\")")
	endpoint := Endpoint{
		Name: "Test",
		Annotations: []annotation.Annotation{
			*ann,
		},
	}
	transport, err := parseHttpTransport("/service", endpoint)

	assert.NotNil(t, err)
	assert.Nil(t, transport)
	assert.Equal(t, "error while parsing endpoint `Test` : method `INVALID` is not valid", err.Error())
}

func TestParseHttpRequestWithValidAnnotations(t *testing.T) {
	ann, _ := annotation.Parse("@http(method=\"GET\", path=\"/test\", request=\"json\", response=\"json\")")
	endpoint := Endpoint{
		Request: &code.Struct{},
		Annotations: []annotation.Annotation{
			*ann,
		},
	}
	request := parseHttpRequest(endpoint)
	assert.NotNil(t, request)
	assert.Equal(t, JSON, request.Format)
}

func TestParseHttpRequestWithNoAnnotations(t *testing.T) {
	endpoint := Endpoint{}
	request := parseHttpRequest(endpoint)
	assert.Nil(t, request)
}

func TestParseHttpRequestWithInvalidAnnotations(t *testing.T) {
	ann, _ := annotation.Parse("@http(method=\"GET\", path=\"/test\", request=\"INVALID\", response=\"json\")")
	endpoint := Endpoint{
		Request: &code.Struct{},
		Annotations: []annotation.Annotation{
			*ann,
		},
	}
	request := parseHttpRequest(endpoint)
	assert.NotNil(t, request)
	assert.Equal(t, JSON, request.Format)
}

func TestHttpresponseWithValidFormat(t *testing.T) {
	format := httpResponseFormat("xml")
	assert.Equal(t, XML, format)
}

func TestHttpresponseWithNoFormat(t *testing.T) {
	format := httpResponseFormat("")
	assert.Equal(t, JSON, format)
}

func TestHttpresponseWithInvalidFormat(t *testing.T) {
	format := httpResponseFormat("INVALID")
	assert.Equal(t, JSON, format)
}

func TestParseHttpRequestParamsWithValidAnnotations(t *testing.T) {
	req := &code.Struct{
		Fields: []code.StructField{
			{
				Parameter: code.Parameter{
					Name: "TestField",
					Type: code.NewType("string"),
				},
				Tags: &code.FieldTags{
					"url": "testUrl",
				},
			},
		},
	}
	request := &HttpRequest{}
	parseHttpRequestParams(req, request)
	assert.True(t, request.HasUrl)
	assert.Equal(t, 1, len(request.Params))
	assert.Equal(t, "TestField", request.Params[0].Field)
	assert.Equal(t, "testUrl", request.Params[0].Name)
	assert.Equal(t, URL, request.Params[0].ParamType)
}

func TestParseHttpRequestParamsWithInvalidAnnotations(t *testing.T) {
	req := &code.Struct{
		Fields: []code.StructField{
			{
				Parameter: code.Parameter{
					Name: "TestField",
					Type: code.NewType("unsupportedType"),
				},
				Tags: &code.FieldTags{
					"url": "testUrl",
				},
			},
		},
	}
	request := &HttpRequest{}
	parseHttpRequestParams(req, request)
	assert.False(t, request.HasUrl)
	assert.Equal(t, 0, len(request.Params))
}

func TestIsQueryTypeSupportedWithSupportedType(t *testing.T) {
	assert.True(t, isQueryTypeSupported("string"))
}

func TestIsQueryTypeSupportedWithUnsupportedType(t *testing.T) {
	assert.False(t, isQueryTypeSupported("unsupportedType"))
}

func TestIsUrlTypeSupportedWithSupportedType(t *testing.T) {
	assert.True(t, isUrlTypeSupported("string"))
}

func TestIsUrlTypeSupportedWithUnsupportedType(t *testing.T) {
	assert.False(t, isUrlTypeSupported("unsupportedType"))
}

func TestParseHttpRequestParamsWithValidQueryAnnotations(t *testing.T) {
	req := &code.Struct{
		Fields: []code.StructField{
			{
				Parameter: code.Parameter{
					Name: "TestField",
					Type: code.NewType("string"),
				},
				Tags: &code.FieldTags{
					"query": "testQuery",
				},
			},
		},
	}
	request := &HttpRequest{}
	parseHttpRequestParams(req, request)
	assert.Equal(t, 1, len(request.Params))
	assert.Equal(t, "TestField", request.Params[0].Field)
	assert.Equal(t, "testQuery", request.Params[0].Name)
	assert.Equal(t, QUERY, request.Params[0].ParamType)
}

func TestParseHttpRequestParamsWithInvalidQueryAnnotations(t *testing.T) {
	req := &code.Struct{
		Fields: []code.StructField{
			{
				Parameter: code.Parameter{
					Name: "TestField",
					Type: code.NewType("unsupportedType"),
				},
				Tags: &code.FieldTags{
					"query": "testQuery",
				},
			},
		},
	}
	request := &HttpRequest{}
	parseHttpRequestParams(req, request)
	assert.Equal(t, 0, len(request.Params))
}

func TestParseHttpRequestParamsWithValidBodyAnnotations(t *testing.T) {
	req := &code.Struct{
		Fields: []code.StructField{
			{
				Parameter: code.Parameter{
					Name: "TestField",
					Type: code.NewType("string"),
				},
				Tags: &code.FieldTags{
					"body": "json",
				},
			},
		},
	}
	request := &HttpRequest{}
	parseHttpRequestParams(req, request)
	assert.True(t, request.HasBody)
	assert.Equal(t, 1, len(request.Params))
	assert.Equal(t, "TestField", request.Params[0].Field)
	assert.Equal(t, "JSON", request.Params[0].Name)
	assert.Equal(t, BODY, request.Params[0].ParamType)
}
