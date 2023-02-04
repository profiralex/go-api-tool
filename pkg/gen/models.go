package gen

import (
	"fmt"
	"github.com/profiralex/go-api-tool/pkg/gen/api"
	"strings"
)

type apiSpec struct {
	Endpoints []api.Endpoint `yaml:"endpoints"`
	Models    []apiModel    `yaml:"models"`
	Module    string        `yaml:"module"`
}

type apiEndpoint struct {
	Path       string                 `yaml:"path"`
	Name       string                 `yaml:"name"`
	Auth       bool                   `yaml:"auth"`
	Method     string                 `yaml:"method"`
	Response   string                 `yaml:"response"`
	BodyFields []apiEndpointBodyField `yaml:"body_fields"`
}

func (e *apiEndpoint) GetURLParams() []string {
	var params []string
	for _, part := range strings.Split(e.Path, "/")[1:] {
		if !strings.Contains(part, "{") {
			continue
		}
		param := strings.ReplaceAll(strings.ReplaceAll(part, "{", ""), "}", "")
		params = append(params, param)
	}
	return params
}

type apiEndpointBodyField struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`
	Validators []string `yaml:"validators"`
}

type apiModel struct {
	Name   string          `yaml:"name"`
	Fields []apiModelField `yaml:"fields"`
}

func (m *apiModel) getField(name string) (apiModelField, bool) {
	for _, field := range m.Fields {
		if field.Name == name {
			return field, true
		}
	}

	return apiModelField{}, false
}

func (m *apiModel) mustGetField(name string) apiModelField {
	field, ok := m.getField(name)
	if !ok {
		panic(fmt.Sprintf("failed to get field %s from model %s", name, m.Name))
	}
	return field
}

func (m *apiModel) hasField(name string) bool {
	_, ok := m.getField(name)
	return ok
}

type apiModelField struct {
	Name           string       `yaml:"name"`
	Type           string       `yaml:"type"`
	Tags           []string     `yaml:"tags"`
	ConstraintsRaw []string     `yaml:"constraints"`
	Constraints    []constraint `yaml:"-"`
}

type constraint struct {
	Name   string
	Value1 string
	Value2 string
	Value3 string
}

func (f *apiModelField) GetGoType() string {
	switch f.Type {
	case "timestamp":
		return "time.Time"
	default:
		return f.Type
	}
}

func (f *apiModelField) parseConstraints() {
	var constraints []constraint
	for _, rawConstraint := range f.ConstraintsRaw {
		c := constraint{}
		parts := strings.Split(rawConstraint, ":")
		c.Name = parts[0]
		if len(parts) > 1 {
			c.Value1 = parts[1]
		}
		if len(parts) > 2 {
			c.Value2 = parts[2]
		}
		if len(parts) > 3 {
			c.Value3 = parts[3]
		}

		constraints = append(constraints, c)
	}

	f.Constraints = constraints
}

func (f *apiModelField) GetConstraint(name string) (constraint, bool) {
	for _, constraint := range f.Constraints {
		if constraint.Name == name {
			return constraint, true
		}
	}

	return constraint{}, false
}

func (f *apiModelField) HasConstraint(name string) bool {
	_, ok := f.GetConstraint(name)
	return ok
}
