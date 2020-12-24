package pkg

import (
	"bytes"
	"fmt"
	"go/format"
	"path"
	"text/template"
)

var serverStatics = map[string]string{
	"server/server.go.tpl":   path.Join(genDirectory, "server", "server.go"),
	"server/response.go.tpl": path.Join(genDirectory, "server", "response.go"),
}

var routesTemplate = template.Must(template.New("routes").
	Funcs(getTemplateFuncs()).
	Parse(templatesBox.MustString("server/routes.go.tpl")))

type serverGenerator struct {
	projectPath string
	apiSpec     apiSpec
}

func NewServerGenerator(projectPath string, spec apiSpec) *serverGenerator {
	return &serverGenerator{
		projectPath: projectPath,
		apiSpec:     spec,
	}
}

func (g *serverGenerator) Generate() error {
	if err := copyBoxStatics(g.projectPath, serverStatics); err != nil {
		return fmt.Errorf("failed to copy static files: %w", err)
	}

	if err := g.generateRoutes(); err != nil {
		return fmt.Errorf("failed to generate routes: %w", err)
	}

	return nil
}

func (g *serverGenerator) generateRoutes() error {
	code := &bytes.Buffer{}

	err := routesTemplate.Execute(code, map[string]interface{}{"apiSpec": g.apiSpec})
	if err != nil {
		return fmt.Errorf("failed to generate routes code: %w", err)
	}

	formattedCodeBytes, err := format.Source(code.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format the generated entity code %w", err)
	}

	routesFilePath := path.Join(g.projectPath, genDirectory, "server", "routes.go")
	return generateFile(routesFilePath, formattedCodeBytes)
}
