package files

import "text/template"

type Provider interface {
	GetFile(filePath string) ([]byte, error)
	MustGetFile(filePath string) []byte
	GetTemplate(filePath string) (*template.Template, error)
	MustGetTemplate(filePath string) *template.Template
}
