CREATE TABLE {{.Table.Name}} (
{{ range $i, $foreignKey := .Table.ForeignKeys -}}
    CONSTRAINT {{ $foreignKey.Field }} FOREIGN KEY ({{ $foreignKey.Field }}) REFERENCES {{ $foreignKey.TableName }} ({{ $foreignKey.ForeignField }}),
{{ end -}}
{{ range $index := .Table.Indexes -}}
    UNIQUE KEY {{ $index.Name }} (
    {{- range $i, $column := $index.Columns -}}
        {{- if gt $i 0 }}, {{ end -}}
        {{$column}}
    {{- end -}}
    ),
{{ end -}}
{{ range $i, $column := .Table.Columns -}}
    {{ $column.Name}} {{$column.Type}}
    {{- if not $column.Nullable -}}
        {{" "}}NOT NULL
    {{- end }}
    {{- if $column.AutoIncrement -}}
        {{" "}}AUTO_INCREMENT
    {{- end }}
    {{- if $column.PrimaryKey -}}
        {{" "}}PRIMARY KEY
    {{- end }}
    {{- if $column.DefaultExpr -}}
        {{" "}}DEFAULT {{ $column.DefaultExpr }}
    {{- end }}
    {{- if $column.OnUpdateExpr -}}
        {{" "}}ON UPDATE {{ $column.OnUpdateExpr }}
    {{- end }}
    {{- if not (last $i $.Table.Columns) -}}
        {{ "," }}
    {{- end }}
{{ end -}}
);