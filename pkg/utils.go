package pkg

import (
	"bytes"
	"fmt"
	"github.com/gertd/go-pluralize"
	"os"
	"strings"
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
