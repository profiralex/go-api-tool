package files

import (
	"bytes"
	"fmt"
	"github.com/profiralex/go-api-tool/pkg/gen/utils"
	"go/format"
	"path/filepath"
	"strings"
)

func CopyFiles(provider Provider, rootDir string, files map[string]string) error {
	for src, destRelativePath := range files {
		filePath := filepath.Join(rootDir, destRelativePath)
		if err := utils.GenerateFile(filePath, provider.MustGetFile(src)); err != nil {
			return fmt.Errorf("failed to copy %s => %s: %w", src, filePath, err)
		}
	}
	return nil
}

func GenerateFiles(provider Provider, rootDir string, files map[string]string, data map[string]interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	for src, destRelativePath := range files {
		filePath := filepath.Join(rootDir, destRelativePath)
		bs, err := ExecuteTemplate(provider, src, data[src])
		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", filePath, err)
		}

		// format only .go.tpl files
		if strings.HasSuffix(filePath, ".go.tpl") {
			formattedCodeBytes, err := format.Source(bs)
			if err != nil {
				return fmt.Errorf("failed to format the generated entity code %w", err)
			}
			bs = formattedCodeBytes
		}

		err = utils.GenerateFile(filePath, bs)
		if err != nil {
			return fmt.Errorf("failed to generate %s : %w", filePath, err)
		}
	}

	return nil
}

func ExecuteTemplate(provider Provider, template string, data interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	tpl := provider.MustGetTemplate(template)
	err := tpl.Execute(buf, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template %s: %w", template, err)
	}

	return buf.Bytes(), nil
}
