package pkg

import (
	"fmt"
	"os"
	"path"
)

const ApiSpecFilename = "api.yml"

var genDirectory = "gen"

type Generator struct {
	projectPath string
	spec        apiSpec
}

func NewGenerator(projectPath string) *Generator {
	return &Generator{
		projectPath: projectPath,
	}
}

func (g *Generator) Init() error {
	apiSpecFilepath := path.Join(g.projectPath, ApiSpecFilename)
	fileInfo, err := os.Stat(apiSpecFilepath)
	if err != nil || fileInfo.IsDir() {
		return fmt.Errorf("provided path is not a project dir")
	}

	f, err := os.Open(apiSpecFilepath)
	if err != nil {
		return fmt.Errorf("failed to open api spec file %w", err)
	}

	parser := &yamlSpecParser{}
	g.spec, err = parser.Parse(f)
	if err != nil {
		return fmt.Errorf("failed to parse api spec file %w", err)
	}

	return nil
}

func (g *Generator) Generate() error {
	err := ensureDirectoryExists(path.Join(g.projectPath, genDirectory))
	if err != nil {
		return fmt.Errorf("failed to create gen directory: %w", err)
	}

	err = g.GenerateMigrations()
	if err != nil {
		return fmt.Errorf("failed to generate migrations: %w", err)
	}

	err = g.GenerateEntities()
	if err != nil {
		return fmt.Errorf("failed to generate entities: %w", err)
	}

	return nil
}

func (g *Generator) GenerateMigrations() error {
	m := NewMigrationsGenerator(g.projectPath, g.spec)
	return m.Generate()
}

func (g *Generator) GenerateEntities() error {
	m := NewEntitiesGenerator(g.projectPath, g.spec)
	return m.Generate()
}
