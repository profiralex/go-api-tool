package pkg

import "strings"

type apiSpec struct {
	Endpoints []apiEndpoint `yaml:"endpoints"`
	Models    []apiModel    `yaml:"models"`
}

type apiEndpoint struct {
	Path       string                 `yaml:"path"`
	Method     string                 `yaml:"method"`
	Response   string                 `yaml:"response"`
	BodyFields []apiEndpointBodyField `yaml:"body_fields"`
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

func (f *apiModelField) getConstraint(name string) (constraint, bool) {
	for _, constraint := range f.Constraints {
		if constraint.Name == name {
			return constraint, true
		}
	}

	return constraint{}, false
}
