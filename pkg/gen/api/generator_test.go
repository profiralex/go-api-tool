package api

import (
	"github.com/profiralex/go-api-tool/pkg/files"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"testing"
)

var log = logrus.New()

func TestGeneratorReturn(t *testing.T) {
	dir, _ := os.MkdirTemp("", "")
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Errorf("failed to cleanup test directory %s : %s", path, err.Error())
		}
	}(dir)

	err := NewGenerator(files.NewBoxProvider()).Generate(dir, Spec{})

	if err != nil {
		t.Errorf("failed to generate files: %s", err)
	}
}

func TestGeneratorGeneratedFiles(t *testing.T) {
	dir, _ := os.MkdirTemp("", "")
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Errorf("failed to cleanup test directory %s : %s", path, err.Error())
		}
	}(dir)

	_ = NewGenerator(files.NewBoxProvider()).Generate(dir, Spec{})

	fileInfos, _ := os.ReadDir(filepath.Join(dir, "server"))
	if len(fileInfos) != 3 {
		t.Errorf("wrong number of generated files %d", len(fileInfos))
	}
}
