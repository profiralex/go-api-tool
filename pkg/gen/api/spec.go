package api

import (
	"strings"
)

type Spec struct {
	Module    string     `yaml:"module"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	Path       string      `yaml:"path"`
	Name       string      `yaml:"name"`
	Auth       bool        `yaml:"auth"`
	Method     string      `yaml:"method"`
	Response   string      `yaml:"response"`
	BodyFields []BodyField `yaml:"body_fields"`
}

func (e *Endpoint) GetURLParams() []string {
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

type BodyField struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`
	Validators []string `yaml:"validators"`
}
