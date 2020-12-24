package pkg

import (
	"bytes"
	"fmt"
	"github.com/gertd/go-pluralize"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
	"unicode"
)

var pluralizeClient = pluralize.NewClient()

func toSnakeCase(s string) string {
	var output strings.Builder
	for i, r := range bytes.Runes([]byte(s)) {
		if i != 0 && unicode.IsUpper(r) {
			_, _ = output.WriteRune('_')
		}
		output.WriteRune(r)
	}

	return strings.ToLower(output.String())
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, _ := range parts {
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}

func toCamelCase(s string) string {
	if strings.Contains(s, "_") {
		s = toPascalCase(s)
	}
	return strings.ToLower(s[0:1]) + s[1:]
}

func plural(s string) string {
	return pluralizeClient.Plural(s)
}

func singular(s string) string {
	return pluralizeClient.Singular(s)
}

func getModelSqlTable(model string) string {
	return toSnakeCase(plural(model))
}

func getModelNameFromSqlTable(model string) string {
	return singular(toPascalCase(model))
}

func getFieldSqlColumn(field string) string {
	return toSnakeCase(field)
}

func getFieldNameFromSqlColumn(sqlColumn string) string {
	return toPascalCase(sqlColumn)
}

func ensureDirectoryExists(path string) error {
	fileStat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to check if migrations directory exists: %w", err)
		}
	} else {
		if !fileStat.IsDir() {
			// oops looks like there is a migrations file
			return fmt.Errorf("migrations file found in project root")
		}
		// directory already exists all good
		return nil
	}

	err = os.Mkdir(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func cleanupDirectory(path string) error {
	ok, err := pathExists(path)
	if err != nil {
		return fmt.Errorf("failed to check if path exists: %w", err)
	}

	if !ok {
		return nil
	}

	err = os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("failed to remove all from directory")
	}

	return nil
}

func generateFile(filePath string, content []byte) error {
	dirPath, _ := path.Split(filePath)

	err := ensureDirectoryExists(dirPath)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	err = ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s file: %w", filePath, err)
	}

	return nil
}

func copyBoxStatics(rootDir string, files map[string]string) error {
	for srcBox, destRelativePath := range files {
		filePath := path.Join(rootDir, destRelativePath)
		if err := generateFile(filePath, templatesBox.MustBytes(srcBox)); err != nil {
			return fmt.Errorf("failed to copy %s => %s: %w", srcBox, filePath, err)
		}
	}
	return nil
}

func getTemplateFuncs(addFuncs ...template.FuncMap) template.FuncMap {
	funcs := template.FuncMap{
		"plural":       plural,
		"singular":     singular,
		"toCamelCase":  toCamelCase,
		"toPascalCase": toPascalCase,
		"toSnakeCase":  toSnakeCase,
		"join":         func(ss []string) string { return strings.Join(ss, ",") },
	}

	if len(addFuncs) > 0 {
		for _, addFuncsMap := range addFuncs {
			for k, v := range addFuncsMap {
				funcs[k] = v
			}
		}
	}

	return funcs
}
