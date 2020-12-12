package pkg

import (
	"bytes"
	"fmt"
	"go/format"
	"path"
	"strings"
	"text/template"
)

var serverContent = `/*Generated code do not modify it*/
package server

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type Server struct {
	server   *http.Server
	serverWG *sync.WaitGroup
}

func NewServer() *Server {
	return &Server{
		server:   &http.Server{},
		serverWG: &sync.WaitGroup{},
	}
}

func (s *Server) SetHandler(router http.Handler) *Server {
	s.server.Handler = router
	return s
}

func (s *Server) SetPort(port int64) *Server {
	s.server.Addr = fmt.Sprintf(":%d", port)
	return s
}

func (s *Server) Start() {
	go func() {
		s.serverWG.Add(1)
		defer s.serverWG.Done()

		log.Infof("Starting server at address %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Errorf("Server failed: %w", err)
		}
	}()
}

func (s *Server) Stop(wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}

	log.Infof("stopping the http server")
	err := s.server.Shutdown(context.Background())
	if err != nil {
		log.Errorf("failed to shutdown the http server: %s", err)
	}

	s.serverWG.Wait()
}
`

var responseContent = `/*Generated code do not modify it*/
package server

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Response the struct to hold the common response field
type Response struct {
	Data   interface{} ` + "`json:\"data\"`" + `
	Errors []APIError  ` + "`json:\"errors\"`" + `
	Status int         ` + "`json:\"status\"`" + `
}

// APIError holds response error information
type APIError struct {
	Message   string ` + "`json:\"message\"`" + `
	Field     string ` + "`json:\"field\"`" + `
	Reference string ` + "`json:\"ref\"`" + `
}

func (response *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, response.Status)
	if response.Status == 500 {
		for _, e := range response.Errors {
			log.Warnf("API ERROR: %s %s %s", e.Reference, e.Field, e.Message)
		}
	}
	return nil
}

func CreateSuccessResponse(data interface{}, status ...int) Response {
	finalStatus := http.StatusOK
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return Response{
		Data:   data,
		Errors: nil,
		Status: finalStatus,
	}
}

func CreateAPIErrorsResponse(errors []APIError, status ...int) Response {
	finalStatus := http.StatusInternalServerError
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return Response{
		Data:   nil,
		Errors: errors,
		Status: finalStatus,
	}
}

func CreateAPIErrorResponse(err APIError, status ...int) Response {
	finalStatus := http.StatusInternalServerError
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return Response{
		Data:   nil,
		Errors: []APIError{err},
		Status: finalStatus,
	}
}

func CreateErrorResponse(err error, status ...int) Response {
	finalStatus := http.StatusInternalServerError
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return Response{
		Data:   nil,
		Errors: []APIError{{Message: err.Error()}},
		Status: finalStatus,
	}
}

func RespondSuccess(w http.ResponseWriter, r *http.Request, data interface{}, status ...int) {
	rsp := CreateSuccessResponse(data, status...)
	_ = render.Render(w, r, &rsp)
}

func RespondAPIError(w http.ResponseWriter, r *http.Request, err APIError, status ...int) {
	rsp := CreateAPIErrorResponse(err, status...)
	_ = render.Render(w, r, &rsp)
}

func RespondError(w http.ResponseWriter, r *http.Request, err error, status ...int) {
	rsp := CreateErrorResponse(err, status...)
	_ = render.Render(w, r, &rsp)
}

func RespondValidationError(w http.ResponseWriter, r *http.Request, err error, status ...int) {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		RespondError(w, r, err, status...)
		return
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		RespondError(w, r, err, status...)
		return
	}

	var apiErrors []APIError
	for _, fieldError := range validationErrors {

		apiError := APIError{
			Field:   fieldError.Field(),
			Message: fieldError.Error(),
		}
		apiErrors = append(apiErrors, apiError)
	}

	rsp := CreateAPIErrorsResponse(apiErrors, status...)
	_ = render.Render(w, r, &rsp)
}
`

var routesTemplate = template.Must(template.New("routes").
	Funcs(template.FuncMap{
		"plural":       plural,
		"singular":     singular,
		"toCamelCase":  toCamelCase,
		"toPascalCase": toPascalCase,
		"toSnakeCase":  toSnakeCase,
		"join":         func(ss []string) string { return strings.Join(ss, ",") },
	}).Parse(`/*Generated code do not modify it*/
package server

import (
	"{{.apiSpec.Module}}/gen/entities"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/go-chi/render"
	"net/http"
)

var validate = validator.New()

type controller interface {
{{- range $i, $endpoint := .apiSpec.Endpoints}} 
	{{$endpoint.Name}}(w http.ResponseWriter, r *http.Request
	{{- range $i, $param := $endpoint.GetURLParams}}, {{$param}} entities.{{toPascalCase $param}} {{- end}}
	{{- if $endpoint.BodyFields}}, data {{$endpoint.Name}}Request{{- end}})
{{- end}}
}

type modelResolver interface {
{{- range $i, $model := .apiSpec.Models}} 
	Find{{$model.Name}}ByUuid(ctx context.Context, uuid string) (entities.{{$model.Name}}, bool, error)
{{- end}}
}

type Handler struct {
	Controller controller
	ModelResolver modelResolver
}

func (h* Handler) RegisterUnauthorizedRoutes(r chi.Router) {
{{- range $i, $endpoint := .apiSpec.Endpoints}} 
	{{- if not $endpoint.Auth}}
	r.{{$endpoint.Method}}("{{$endpoint.Path}}", h.{{$endpoint.Name}})
	{{- end}}
{{- end}}
}

func (h* Handler) RegisterAuthorizedRoutes(r chi.Router) {
{{- range $i, $endpoint := .apiSpec.Endpoints}} 
	{{- if $endpoint.Auth}}
	r.{{$endpoint.Method}}("{{$endpoint.Path}}", h.{{$endpoint.Name}})
	{{- end}}
{{- end}}
}

{{- range $i, $endpoint := .apiSpec.Endpoints}} 
	{{if $endpoint.BodyFields}}
		type {{$endpoint.Name}}Request struct {
		{{- range $i, $bodyField := $endpoint.BodyFields}} 
			{{toPascalCase $bodyField.Name}} {{$bodyField.Type}} ` + "`json:\"{{$bodyField.Name}}\" validate:\"{{join $bodyField.Validators}}\"`" + `
		{{- end}}
		}
		
		func (d *{{$endpoint.Name}}Request) Bind(*http.Request) error {
			return validate.Struct(d)
		}
	{{- end}}

func (h* Handler) {{$endpoint.Name}}(w http.ResponseWriter, r *http.Request) {
	{{- range $i, $param := $endpoint.GetURLParams}} 
		{{toCamelCase $param}}Uuid := chi.URLParam(r, "{{$param}}")
		{{toCamelCase $param}}, ok, err := h.ModelResolver.Find{{toPascalCase $param}}ByUuid(r.Context(), {{toCamelCase $param}}Uuid)
		if err != nil {
			RespondError(w, r, fmt.Errorf("failed to get {{toCamelCase $param}}: %w", err))
			return
		} else if !ok {{- if gt $i 0}} || {{toCamelCase (index $endpoint.GetURLParams 0)}}.Id != {{toCamelCase $param}}.{{toPascalCase (index $endpoint.GetURLParams 0)}}Id {{- end}} {
			RespondError(w, r, fmt.Errorf("{{toCamelCase $param}} not found"), http.StatusNotFound)
			return
		}
	{{end}}

	{{- if $endpoint.BodyFields}}
		data := {{$endpoint.Name}}Request{}
		if err := render.Bind(r, &data); err != nil {
			RespondValidationError(w, r, err, http.StatusBadRequest)
			return
		}
	{{- end}}

	h.Controller.{{$endpoint.Name}}(w, r
	{{- range $i, $param := $endpoint.GetURLParams}}, {{$param}} {{- end}}
	{{- if $endpoint.BodyFields}}, data {{- end}})
}
{{- end}}
`))

type serverGenerator struct {
	projectPath string
	apiSpec     apiSpec
}

func NewServerGenerator(projectPath string, spec apiSpec) *serverGenerator {
	return &serverGenerator{
		projectPath: projectPath,
		apiSpec:     spec,
	}
}

func (g *serverGenerator) Generate() error {
	if err := g.generateServer(); err != nil {
		return fmt.Errorf("failed to generate http server: %w", err)
	}

	if err := g.generateResponseHelper(); err != nil {
		return fmt.Errorf("failed to generate response helper: %w", err)
	}

	if err := g.generateRoutes(); err != nil {
		return fmt.Errorf("failed to generate routes: %w", err)
	}

	return nil
}

func (g *serverGenerator) generateServer() error {
	serverFilePath := path.Join(g.projectPath, genDirectory, "server", "server.go")
	return generateFile(serverFilePath, []byte(serverContent))
}

func (g *serverGenerator) generateResponseHelper() error {
	responseHelperFilePath := path.Join(g.projectPath, genDirectory, "server", "response.go")
	return generateFile(responseHelperFilePath, []byte(responseContent))
}

func (g *serverGenerator) generateRoutes() error {
	code := &bytes.Buffer{}

	err := routesTemplate.Execute(code, map[string]interface{}{"apiSpec": g.apiSpec})
	if err != nil {
		return fmt.Errorf("failed to generate routes code: %w", err)
	}

	formattedCodeBytes, err := format.Source(code.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format the generated entity code %w", err)
	}

	routesFilePath := path.Join(g.projectPath, genDirectory, "server", "routes.go")
	return generateFile(routesFilePath, formattedCodeBytes)
}
