package utils

import (
	"bytes"
	"fmt"
	"github.com/gertd/go-pluralize"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var pluralizeClient = pluralize.NewClient()

func ToSnakeCase(s string) string {
	var output strings.Builder
	for i, r := range bytes.Runes([]byte(s)) {
		if i != 0 && unicode.IsUpper(r) {
			_, _ = output.WriteRune('_')
		}
		output.WriteRune(r)
	}

	return strings.ToLower(output.String())
}

func ToPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, _ := range parts {
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}

func ToCamelCase(s string) string {
	if strings.Contains(s, "_") {
		s = ToPascalCase(s)
	}
	return strings.ToLower(s[0:1]) + s[1:]
}

func Plural(s string) string {
	return pluralizeClient.Plural(s)
}

func Singular(s string) string {
	return pluralizeClient.Singular(s)
}

func GetModelSqlTable(model string) string {
	return ToSnakeCase(Plural(model))
}

func GetModelNameFromSqlTable(model string) string {
	return Singular(ToPascalCase(model))
}

func GetFieldSqlColumn(field string) string {
	return ToSnakeCase(field)
}

func GetFieldNameFromSqlColumn(sqlColumn string) string {
	return ToPascalCase(sqlColumn)
}

func EnsureDirectoryExists(path string) error {
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

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func CleanupDirectory(path string) error {
	ok, err := PathExists(path)
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

func GenerateFile(filePath string, content []byte) error {
	dirPath, _ := filepath.Split(filePath)

	err := EnsureDirectoryExists(dirPath)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s file: %w", filePath, err)
	}

	return nil
}
