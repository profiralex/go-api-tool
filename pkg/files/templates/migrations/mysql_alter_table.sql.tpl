{{if .Diff.HasAddedColumns -}}
    ALTER TABLE {{.Table.Name}} (
    {{ range $i, $column := .Diff.AddedColumns -}}
        ADD COLUMN {{ $column.Name}} {{$column.Type}}
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
{{- end}}
{{if .Diff.HasUpdatedColumns -}}
    ALTER TABLE {{.Table.Name}} (
    {{ range $i, $column := .Diff.AddedColumns -}}
        ALTER COLUMN {{ $column.Name}} {{$column.Type}}
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
{{- end}}
