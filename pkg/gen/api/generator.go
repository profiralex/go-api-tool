package api

import (
	"fmt"
	"github.com/profiralex/go-api-tool/pkg/files"
	"path/filepath"
)

var serverStatics = map[string]string{
	filepath.Join("templates", "server", "server.go.tpl"):   filepath.Join("server", "server.go"),
	filepath.Join("templates", "server", "response.go.tpl"): filepath.Join("server", "response.go"),
}

type Generator struct {
	fp files.Provider
}

func NewGenerator(fp files.Provider) *Generator {
	return &Generator{
		fp: fp,
	}
}

func (g *Generator) Generate(genPath string, spec Spec) error {
	if err := files.CopyFiles(g.fp, genPath, serverStatics); err != nil {
		return fmt.Errorf("failed to copy static files: %w", err)
	}

	routesPath := filepath.Join("templates", "server", "routes.go.tpl")

	filesToGenerate := map[string]string{
		routesPath: filepath.Join("server", "routes.go"),
	}
	templatesData := map[string]interface{}{
		routesPath: map[string]interface{}{"Spec": spec},
	}
	if err := files.GenerateFiles(g.fp, genPath, filesToGenerate, templatesData); err != nil {
		return fmt.Errorf("failed to generate routes: %w", err)
	}

	return nil
}
