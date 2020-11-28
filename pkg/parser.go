package pkg

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

type specReader interface {
	io.Reader
}

type specParser interface {
	Parse(r specReader) (apiSpec, error)
}

// ensure yamlSpecParser implements the interface
var _ specParser = &yamlSpecParser{}

type yamlSpecParser struct {
}

func (y *yamlSpecParser) Parse(r specReader) (apiSpec, error) {
	var result apiSpec
	err := yaml.NewDecoder(r).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("failed to parse api spec: %w", err)
	}

	for modelIndex, model := range result.Models {
		for fieldIndex, field := range model.Fields {
			field.parseConstraints()
			result.Models[modelIndex].Fields[fieldIndex] = field
		}
	}

	return result, nil
}
