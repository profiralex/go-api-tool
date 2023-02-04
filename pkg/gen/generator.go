package gen

import (
	"fmt"
	"github.com/profiralex/go-api-tool/pkg/gen/api"
	"github.com/profiralex/go-api-tool/pkg/gen/entities"
	"github.com/profiralex/go-api-tool/pkg/gen/logs"
	"github.com/profiralex/go-api-tool/pkg/gen/migrations"
	"github.com/profiralex/go-api-tool/pkg/gen/utils"
	"os"
	"path/filepath"
)

const ApiSpecFilename = "api.yml"

var genDirectory = "gen"

type Generator struct {
	apiGenerator  api.Generator
	logsGenerator logs.Generator
}

func NewGenerator(projectPath string) *Generator {
	return &Generator{}
}

func (g *Generator) Init() error {
	return nil
}

func (g *Generator) Generate(projectPath string) error {
	apiSpecFilepath := filepath.Join(projectPath, ApiSpecFilename)
	fileInfo, err := os.Stat(apiSpecFilepath)
	if err != nil || fileInfo.IsDir() {
		return fmt.Errorf("provided path is not a project dir")
	}

	f, err := os.Open(apiSpecFilepath)
	if err != nil {
		return fmt.Errorf("failed to open api spec file %w", err)
	}

	parser := &yamlSpecParser{}
	spec, err = parser.Parse(f)
	if err != nil {
		return fmt.Errorf("failed to parse api spec file %w", err)
	}

	err := g.prepareGenDirectory()
	if err != nil {
		return fmt.Errorf("failed to prepare gen directory: %w", err)
	}

	err = g.GenerateMigrations()
	if err != nil {
		return fmt.Errorf("failed to generate migrations: %w", err)
	}

	err = g.GenerateEntities()
	if err != nil {
		return fmt.Errorf("failed to generate entities: %w", err)
	}

	err = g.GenerateLogs()
	if err != nil {
		return fmt.Errorf("failed to generate logs: %w", err)
	}

	err = g.GenerateServer()
	if err != nil {
		return fmt.Errorf("failed to generate server: %w", err)
	}

	return nil
}

func (g *Generator) prepareGenDirectory() error {
	genPath := filepath.Join(g.projectPath, genDirectory)
	err := utils.CleanupDirectory(genPath)
	if err != nil {
		return fmt.Errorf("failed to cleanup the gen directory: %w", err)
	}

	err = utils.EnsureDirectoryExists(genPath)
	if err != nil {
		return fmt.Errorf("failed to create the gen directory: %w", err)
	}

	return nil
}

func (g *Generator) GenerateMigrations() error {
	m := migrations.NewMigrationsGenerator(g.projectPath, g.spec)
	return m.Generate()
}

func (g *Generator) GenerateEntities() error {
	m := entities.NewEntitiesGenerator(g.projectPath, g.spec)
	return m.Generate()
}

func (g *Generator) GenerateLogs() error {
	m := logs.NewGenerator(g.projectPath, g.spec)
	return m.Generate()
}

func (g *Generator) GenerateServer() error {
	m := api.NewGenerator(g.projectPath, g.spec)
	return m.Generate()
}
