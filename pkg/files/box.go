package files

import (
	"embed"
	"fmt"
	"github.com/profiralex/go-api-tool/pkg/gen/utils"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"text/template"
)

//go:embed templates/*
var templatesFS embed.FS

var templateFuncMap = template.FuncMap{
	"plural":       utils.Plural,
	"singular":     utils.Singular,
	"toCamelCase":  utils.ToCamelCase,
	"toPascalCase": utils.ToPascalCase,
	"toSnakeCase":  utils.ToSnakeCase,
	"join":         func(ss []string) string { return strings.Join(ss, ",") },
	"last":         func(x int, a interface{}) bool { return x == reflect.ValueOf(a).Len()-1 },
}

type EmbedProvider struct {
	cache map[string]*template.Template
	lock  sync.RWMutex
}

func NewBoxProvider() *EmbedProvider {
	return &EmbedProvider{
		cache: map[string]*template.Template{},
	}
}

func (b *EmbedProvider) GetFile(filePath string) ([]byte, error) {
	return templatesFS.ReadFile(filepath.ToSlash(filePath))
}

func (b *EmbedProvider) MustGetFile(filePath string) []byte {
	bs, err := b.GetFile(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to get file %s : %w", filePath, err))
	}
	return bs
}

func (b *EmbedProvider) GetTemplate(filePath string) (*template.Template, error) {
	b.lock.RLock()
	tpl, ok := b.cache[filePath]
	b.lock.RUnlock()
	if ok {
		return tpl, nil
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	// additional check just in case
	tpl, ok = b.cache[filePath]
	if ok {
		return tpl, nil
	}

	content, err := b.GetFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content %w", err)
	}

	tpl, err = template.New(filePath).
		Funcs(templateFuncMap).
		Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	b.cache[filePath] = tpl
	return tpl, nil
}

func (b *EmbedProvider) MustGetTemplate(filePath string) *template.Template {
	tpl, err := b.GetTemplate(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to get template %s : %w", filePath, err))
	}

	return tpl
}
