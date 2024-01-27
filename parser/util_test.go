package parser

import (
	"github.com/go-services/annotation"
	"github.com/go-services/code"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindAnnotations_WithMatchingName(t *testing.T) {
	annotations := []annotation.Annotation{
		{Name: "match"},
		{Name: "no_match"},
	}
	found := findAnnotations("match", annotations)

	assert.Len(t, found, 1, "should be equal")
	assert.Equal(t, "match", found[0].Name, "should be equal")
}

func TestFindAnnotations_WithNoMatchingName(t *testing.T) {
	annotations := []annotation.Annotation{
		{Name: "no_match"},
	}
	found := findAnnotations("match", annotations)

	assert.Len(t, found, 0, "should be equal")
}

func TestGetTag_WithExistingKey(t *testing.T) {
	tags := code.FieldTags{
		"key": "value",
	}
	tag := getTag("key", tags)

	assert.Equal(t, "value", tag, "should be equal")
}

func TestGetTag_WithNonExistingKey(t *testing.T) {
	tags := code.FieldTags{
		"key": "value",
	}
	tag := getTag("non_existing_key", tags)

	assert.Equal(t, "", tag, "should be equal")
}

func TestGetParameter_WithRequired(t *testing.T) {
	name, required := getParameter("param,required")

	assert.Equal(t, "param", name, "should be equal")
	assert.True(t, required, "should be true")
}

func TestGetParameter_WithoutRequired(t *testing.T) {
	name, required := getParameter("param")

	assert.Equal(t, "param", name, "should be equal")
	assert.False(t, required, "should be false")
}

func TestGetParameter_WithEmptyString(t *testing.T) {
	name, required := getParameter("")

	assert.Equal(t, "", name, "should be equal")
	assert.False(t, required, "should be false")
}
