package pkg

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"path"
	"strings"
)

var genDirectory = "gen"

type modelsGenerator struct {
	projectPath string
	apiSpec     apiSpec
}

func NewModelsGenerator(projectPath string, spec apiSpec) *modelsGenerator {
	return &modelsGenerator{
		projectPath: projectPath,
		apiSpec:     spec,
	}
}

func (g *modelsGenerator) Generate() error {
	err := ensureDirectoryExists(path.Join(g.projectPath, genDirectory))
	if err != nil {
		return fmt.Errorf("failed to create gen directory: %w", err)
	}

	content, err := g.generateModelsCode()
	if err != nil {
		return fmt.Errorf("failed to create models file: %w", err)
	}

	filePath := path.Join(g.projectPath, genDirectory, "models.go")
	err = ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write models file: %w", err)
	}

	return nil
}

func (g *modelsGenerator) generateModelsCode() ([]byte, error) {
	fileHeader := "/*Generated code do not modify it*/\npackage gen"
	imports := "import \"time\""
	var models []string

	for _, model := range g.apiSpec.Models {
		var fields []string
		for _, field := range model.Fields {
			field := fmt.Sprintf("%s %s `dbr:\"%s\" json:\"%s\"`", field.Name, g.getGoType(field.Type), getFieldSqlColumn(field.Name), getFieldSqlColumn(field.Name))
			fields = append(fields, field)
		}

		model := fmt.Sprintf(`type %s struct {
%s
}
`, model.Name, strings.Join(fields, "\n"))
		models = append(models, model)
	}

	result, err := format.Source([]byte(fileHeader + "\n" + imports + "\n" + strings.Join(models, "\n")))
	if err != nil {
		return nil, fmt.Errorf("failed to format the code %w", err)
	}

	return result, nil
}

func (g *modelsGenerator) getGoType(fieldType string) string {
	switch fieldType {
	case "timestamp":
		return "time.Time"
	default:
		return fieldType
	}
}
