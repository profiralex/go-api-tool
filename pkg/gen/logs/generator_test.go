package logs

import (
	"github.com/profiralex/go-api-tool/pkg/files"
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratorReturn(t *testing.T) {
	dir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(dir)

	err := NewGenerator(files.NewBoxProvider()).Generate(dir)

	if err != nil {
		t.Errorf("failed to generate files: %s", err)
	}
}

func TestGeneratorGeneratedFiles(t *testing.T) {
	dir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(dir)

	_ = NewGenerator(files.NewBoxProvider()).Generate(dir)

	fileInfos, _ := os.ReadDir(filepath.Join(dir, "logs"))
	if len(fileInfos) != 1 {
		t.Errorf("wrong number of generated files")
	}
}
