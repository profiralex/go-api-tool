/*Generated code do not modify it*/
package entities

import (
"context"
"fmt"
"github.com/gocraft/dbr/v2"
"github.com/google/uuid"
"time"
)

type Base{{plural .Model.Name}}Repo interface {
FindById(ctx context.Context, id int64) ({{$.Model.Name}}, bool, error)
FindByIds(ctx context.Context, ids ...int64) ([]{{$.Model.Name}}, error)
FindByUuid(ctx context.Context, uuid string) ({{$.Model.Name}}, bool, error)
FindByUuids(ctx context.Context, uuids ...string) ([]{{$.Model.Name}}, error)
{{- range $i, $field := .Model.Fields}}
    {{- if $field.HasConstraint "foreign_key"}}
        FindBy{{$field.Name}}(ctx context.Context, {{toCamelCase $field.Name}} {{$field.GetGoType}}) ({{$.Model.Name}}, bool, error)
        FindBy{{$field.Name}}s(ctx context.Context, {{toCamelCase $field.Name}}s ...{{$field.GetGoType}}) ([]{{$.Model.Name}}, error)
    {{- end}}
{{- end}}
Save(ctx context.Context, records ...*{{$.Model.Name}}) error
Delete(ctx context.Context, records ...*{{$.Model.Name}}) error
DeleteByIds(ctx context.Context, ids ...int64) error
}

type {{$.Model.Name}} struct {
{{- range $i, $field := .Model.Fields}}
    {{$field.Name}} {{$field.GetGoType}} `json:"{{toSnakeCase $field.Name}}" dbr:"{{toSnakeCase $field.Name}}"`
{{- end}}
}

var _ Base{{plural .Model.Name}}Repo = &{{$.Model.Name}}EntitiesManager{}
type {{$.Model.Name}}EntitiesManager struct {
table   string
session dbr.Session
}

func New{{$.Model.Name}}EntitiesManager(session dbr.Session) *{{$.Model.Name}}EntitiesManager {
return &{{$.Model.Name}}EntitiesManager{
table:   "{{getModelSqlTable .Model.Name}}",
session: session,
}
}

func (m *{{$.Model.Name}}EntitiesManager) FindById(ctx context.Context, id int64) ({{$.Model.Name}}, bool, error) {
return m.Query().IdEq(id).FindFirst(ctx)
}

func (m *{{$.Model.Name}}EntitiesManager) FindByIds(ctx context.Context, ids ...int64) ([]{{$.Model.Name}}, error) {
return m.Query().IdIn(ids).FindAll(ctx)
}

func (m *{{$.Model.Name}}EntitiesManager) FindByUuid(ctx context.Context, uuid string) ({{$.Model.Name}}, bool, error) {
return m.Query().UuidEq(uuid).FindFirst(ctx)
}

func (m *{{$.Model.Name}}EntitiesManager) FindByUuids(ctx context.Context, uuids ...string) ([]{{$.Model.Name}}, error) {
return m.Query().UuidIn(uuids).FindAll(ctx)
}

{{- range $i, $field := .Model.Fields}}
    {{ if $field.HasConstraint "foreign_key"}}
        func (m *{{$.Model.Name}}EntitiesManager) FindBy{{$field.Name}}(ctx context.Context, {{toCamelCase $field.Name}} {{$field.GetGoType}}) ({{$.Model.Name}}, bool, error) {
        return m.Query().{{$field.Name}}Eq({{toCamelCase $field.Name}}).FindFirst(ctx)
        }

        func (m *{{$.Model.Name}}EntitiesManager) FindBy{{$field.Name}}s(ctx context.Context, {{toCamelCase $field.Name}}s ...{{$field.GetGoType}}) ([]{{$.Model.Name}}, error) {
        return m.Query().{{$field.Name}}In({{toCamelCase $field.Name}}s).FindAll(ctx)
        }
    {{- end}}
{{- end}}

func (m *{{$.Model.Name}}EntitiesManager) Save(ctx context.Context, records ...*{{$.Model.Name}}) error {
var recordsToInsert []*{{$.Model.Name}}
var recordsToUpdate []*{{$.Model.Name}}

for _, record := range records {
if record.Id == 0 || record.CreatedAt.IsZero() || record.UpdatedAt.IsZero() {
recordsToInsert = append(recordsToInsert, record)
} else {
recordsToUpdate = append(recordsToUpdate, record)
}
}

for _, record := range recordsToUpdate {
err := m.singleUpdate(ctx, record)
if err != nil {
return fmt.Errorf("failed to update record %d: %w", record.Id, err)
}
}

for _, record := range recordsToInsert {
err := m.singleInsert(ctx, record)
if err != nil {
return fmt.Errorf("failed to insert record: %w", err)
}
}

return nil
}

func (m *{{$.Model.Name}}EntitiesManager) singleUpdate(ctx context.Context, record *{{$.Model.Name}}) error {
record.UpdatedAt = time.Now()
err := m.session.Update(m.table).
SetMap(map[string]interface{}{
{{- range $i, $field := .Model.Fields}}
    {{- if and (ne $field.Name "Id") (ne $field.Name "Uuid") }}
        "{{toSnakeCase $field.Name}}": record.{{$field.Name}},
    {{- end}}
{{- end}}
}).
Where("id = ?", record.Id).
LoadContext(ctx, record)

return err
}

func (m *{{$.Model.Name}}EntitiesManager) singleInsert(ctx context.Context, record *{{$.Model.Name}}) error {
record.Uuid = uuid.New().String()
record.CreatedAt = time.Now()
record.UpdatedAt = record.CreatedAt

_, err := m.session.InsertInto(m.table).
Columns(
{{- range $i, $field := .Model.Fields}}
    {{- if ne $field.Name "Id" }}
        "{{toSnakeCase $field.Name}}",
    {{- end}}
{{- end}}
).
Record(record).
ExecContext(ctx)

return err
}

func (m *{{$.Model.Name}}EntitiesManager) Delete(ctx context.Context, records ...*{{$.Model.Name}}) error {
if len(records) == 0 {
return nil
}

var ids []int64
for _, record := range records {
ids = append(ids, record.Id)
}

return m.DeleteByIds(ctx, ids...)
}

func (m *{{$.Model.Name}}EntitiesManager) DeleteByIds(ctx context.Context, ids ...int64) error {
_, err := m.session.DeleteFrom(m.table).
Where(dbr.Eq("id", ids)).
ExecContext(ctx)

if err != nil {
return fmt.Errorf("failed to delete records: %w", err)
}

return nil
}

type {{$.Model.Name}}EntitiesQuery struct {
filters      []dbr.Builder
offset       *uint64
limit        *uint64
orderingAsc  []string
orderingDesc []string
table        string
session      dbr.Session
}

func (m *{{$.Model.Name}}EntitiesManager) Query() *{{$.Model.Name}}EntitiesQuery {
return &{{$.Model.Name}}EntitiesQuery{
table:   m.table,
session: m.session,
}
}

func (q *{{$.Model.Name}}EntitiesQuery) Limit(limit uint64) *{{$.Model.Name}}EntitiesQuery {
q.limit = &limit
return q
}

func (q *{{$.Model.Name}}EntitiesQuery) Offset(offset uint64) *{{$.Model.Name}}EntitiesQuery {
q.offset = &offset
return q
}


func (q *{{$.Model.Name}}EntitiesQuery) FindAll(ctx context.Context) ([]{{$.Model.Name}}, error) {
var records []{{$.Model.Name}}
_, err := q.computeSelectStatement().
LoadContext(ctx, &records)

if err != nil {
return records, fmt.Errorf("failed to load records: %w", err)
}

return records, nil
}

func (q *{{$.Model.Name}}EntitiesQuery) FindFirst(ctx context.Context) ({{$.Model.Name}}, bool, error) {
records, err := q.Limit(1).FindAll(ctx)
if err != nil {
return {{$.Model.Name}}{}, false, fmt.Errorf("failed to get records: %w", err)
}

if len(records) == 0 {
return {{$.Model.Name}}{}, false, nil
}

return records[0], true, nil
}

func (q *{{$.Model.Name}}EntitiesQuery) Exists(ctx context.Context) (bool, error) {
_, ok, err := q.FindFirst(ctx)
if err != nil {
return false, fmt.Errorf("failed to get first record: %w", err)
}

return ok, nil
}

func (q *{{$.Model.Name}}EntitiesQuery) Count(ctx context.Context) (uint64, error) {
result := struct {
Count uint64 `dbr:"count"`
}{}

err := q.computeSelectStatement("COUNT(*) as count").LoadOneContext(ctx, &result)
if err != nil {
return 0, fmt.Errorf("failed to retrieve records count: %w", err)
}

return result.Count, nil
}

func (q *{{$.Model.Name}}EntitiesQuery) computeSelectStatement(columns ...string) *dbr.SelectStmt {
if len(columns) == 0 {
columns = append(columns, "*")
}

stmt := q.session.Select(columns...).From(q.table)

if len(q.filters) == 1 {
stmt = stmt.Where(q.filters[0])
}

if len(q.filters) > 1 {
stmt = stmt.Where(dbr.And(q.filters...))
}

if q.limit != nil {
stmt = stmt.Limit(*q.limit)
}

if q.offset != nil {
stmt = stmt.Offset(*q.offset)
}

for _, column := range q.orderingAsc {
stmt = stmt.OrderAsc(column)
}

for _, column := range q.orderingDesc {
stmt = stmt.OrderDesc(column)
}

return stmt
}

{{- range $i, $field := .Model.Fields}}
    func (q *{{$.Model.Name}}EntitiesQuery) OrderAscBy{{$field.Name}}() *{{$.Model.Name}}EntitiesQuery {
    q.orderingAsc = append(q.orderingAsc, "{{toSnakeCase $field.Name}}")
    return q
    }

    func (q *{{$.Model.Name}}EntitiesQuery) OrderDescBy{{$field.Name}}() *{{$.Model.Name}}EntitiesQuery {
    q.orderingAsc = append(q.orderingDesc, "{{toSnakeCase $field.Name}}")
    return q
    }

    func (q *{{$.Model.Name}}EntitiesQuery) {{$field.Name}}Eq(value {{$field.GetGoType}}) *{{$.Model.Name}}EntitiesQuery {
    q.filters = append(q.filters, dbr.Eq("{{toSnakeCase $field.Name}}", value))
    return q
    }

    func (q *{{$.Model.Name}}EntitiesQuery) {{$field.Name}}Neq(value {{$field.GetGoType}}) *{{$.Model.Name}}EntitiesQuery {
    q.filters = append(q.filters, dbr.Neq("{{toSnakeCase $field.Name}}", value))
    return q
    }

    func (q *{{$.Model.Name}}EntitiesQuery) {{$field.Name}}In(values []{{$field.GetGoType}}) *{{$.Model.Name}}EntitiesQuery {
    q.filters = append(q.filters, dbr.Eq("{{toSnakeCase $field.Name}}", values))
    return q
    }

    func (q *{{$.Model.Name}}EntitiesQuery) {{$field.Name}}Gt(value {{$field.GetGoType}}) *{{$.Model.Name}}EntitiesQuery {
    q.filters = append(q.filters, dbr.Gt("{{toSnakeCase $field.Name}}", value))
    return q
    }

    func (q *{{$.Model.Name}}EntitiesQuery) {{$field.Name}}Gte(value {{$field.GetGoType}}) *{{$.Model.Name}}EntitiesQuery {
    q.filters = append(q.filters, dbr.Gte("{{toSnakeCase $field.Name}}", value))
    return q
    }

    func (q *{{$.Model.Name}}EntitiesQuery) {{$field.Name}}Lt(value {{$field.GetGoType}}) *{{$.Model.Name}}EntitiesQuery {
    q.filters = append(q.filters, dbr.Lt("{{toSnakeCase $field.Name}}", value))
    return q
    }

    func (q *{{$.Model.Name}}EntitiesQuery) {{$field.Name}}Lte(value {{$field.GetGoType}}) *{{$.Model.Name}}EntitiesQuery {
    q.filters = append(q.filters, dbr.Lte("{{toSnakeCase $field.Name}}", value))
    return q
    }

    {{- if eq $field.Type "string" }}
        func (q *{{$.Model.Name}}EntitiesQuery) {{$field.Name}}Like(value string) *{{$.Model.Name}}EntitiesQuery {
        q.filters = append(q.filters, dbr.Like("{{toSnakeCase $field.Name}}", value))
        return q
        }

        func (q *{{$.Model.Name}}EntitiesQuery) {{$field.Name}}NotLike(value string) *{{$.Model.Name}}EntitiesQuery {
        q.filters = append(q.filters, dbr.NotLike("{{toSnakeCase $field.Name}}", value))
        return q
        }
    {{- end}}
{{- end}}
