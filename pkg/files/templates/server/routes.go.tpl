/*Generated code do not modify it*/
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
            {{toPascalCase $bodyField.Name}} {{$bodyField.Type}} `json:"{{$bodyField.Name}}" validate:"{{join $bodyField.Validators}}"`
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
