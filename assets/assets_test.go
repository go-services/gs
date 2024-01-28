package assets

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gs/fs"
	"testing"
)

func TestGetHttpResponseEncodeFunction(t *testing.T) {
	assert.Equal(t, "JsonEncoder", GetHttpResponseEncodeFunction("JSON"))
	assert.Equal(t, "XmlEncoder", GetHttpResponseEncodeFunction("XML"))
	assert.Equal(t, "", GetHttpResponseEncodeFunction("INVALID"))
}

func TestGetHttpRequestDecoderFunction(t *testing.T) {
	assert.Equal(t, "JsonDecoder", GetHttpRequestDecoderFunction("JSON"))
	assert.Equal(t, "XmlDecoder", GetHttpRequestDecoderFunction("XML"))
	assert.Equal(t, "FormDecoder", GetHttpRequestDecoderFunction("FORM"))
	assert.Equal(t, "", GetHttpRequestDecoderFunction("INVALID"))
}

func TestReadTemplate(t *testing.T) {
	_, err := ReadTemplate("nonexistent")
	assert.Error(t, err)

	content, err := ReadTemplate("utils/utils.go.tmpl")
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestParseTemplate(t *testing.T) {
	_, err := ParseTemplate("nonexistent", nil)
	assert.Error(t, err)

	content, err := ParseTemplate("utils/utils.go.tmpl", nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestParseAndWriteTemplate(t *testing.T) {
	fs.SetTestFs(afero.NewMemMapFs())
	err := ParseAndWriteTemplate("nonexistent", "output", nil)
	assert.Error(t, err)

	err = ParseAndWriteTemplate("utils/utils.go.tmpl", "output", nil)
	assert.NoError(t, err)
}
