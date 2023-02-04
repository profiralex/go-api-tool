package entities

import (
	"bytes"
	"fmt"
	"github.com/profiralex/go-api-tool/pkg/gen"
	"go/format"
	"io/ioutil"
	"path/filepath"
)

//
//var entitiesDirectory = "entities"

// var entityTemplate = template.Must(template.New("entity").
//
//	Funcs(gen.getTemplateFuncs(template.FuncMap{
//		"getModelSqlTable": gen.getModelSqlTable,
//	})).
//	Parse(gen.templatesBox.MustString("entities/entity.go.tpl")))
type entitiesGenerator struct {
	projectPath string
	apiSpec     gen.apiSpec
}

func NewEntitiesGenerator(projectPath string, spec gen.apiSpec) *entitiesGenerator {
	return &entitiesGenerator{
		projectPath: projectPath,
		apiSpec:     spec,
	}
}

func (g *entitiesGenerator) Generate() error {
	err := gen.ensureDirectoryExists(filepath.Join(g.projectPath, gen.genDirectory, entitiesDirectory))
	if err != nil {
		return fmt.Errorf("failed to create entitites directory: %w", err)
	}

	for _, model := range g.apiSpec.Models {
		err := g.generateEntityFile(model)
		if err != nil {
			return fmt.Errorf("failed to generate entity file for %s model: %w", model.Name, err)
		}
	}

	return nil
}

func (g *entitiesGenerator) generateEntityFile(model gen.apiModel) error {
	code := &bytes.Buffer{}

	err := entityTemplate.Execute(code, map[string]interface{}{"Table": model})
	if err != nil {
		return fmt.Errorf("failed to generate entity code: %w", err)
	}

	formattedCodeBytes, err := format.Source(code.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format the generated entity code %w", err)
	}

	entityFileName := fmt.Sprintf("%s_entity.go", gen.toSnakeCase(model.Name))
	filePath := filepath.Join(g.projectPath, gen.genDirectory, entitiesDirectory, entityFileName)
	err = ioutil.WriteFile(filePath, formattedCodeBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write models file: %w", err)
	}

	return nil
}

func (g *entitiesGenerator) getGoType(fieldType string) string {
	switch fieldType {
	case "timestamp":
		return "time.Time"
	default:
		return fieldType
	}
}
