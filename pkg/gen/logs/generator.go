package logs

import (
	"github.com/profiralex/go-api-tool/pkg/files"
	"path/filepath"
)

type Generator struct {
	genPath       string
	filesProvider files.Provider
}

var logsStatics = map[string]string{
	filepath.Join("templates", "logs", "logs.go.tpl"): filepath.Join("logs", "logs.go"),
}

func NewGenerator(filesProvider files.Provider) *Generator {
	return &Generator{
		filesProvider: filesProvider,
	}
}

func (g *Generator) Generate(genPath string) error {
	return files.CopyFiles(g.filesProvider, genPath, logsStatics)
}
