package pkg

import (
	"bytes"
	"github.com/gertd/go-pluralize"
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

func plural(s string) string {
	return pluralizeClient.Plural(s)
}

func singular(s string) string {
	return pluralizeClient.Singular(s)
}

func getModelSqlTable(model string) string {
	return toSnakeCase(plural(model))
}

func getFieldSqlColumn(field string) string {
	return toSnakeCase(field)
}
