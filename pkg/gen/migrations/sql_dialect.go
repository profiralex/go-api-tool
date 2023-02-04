package migrations

type SqlDialect interface {
	CreateSql(model Table) ([]byte, error)
	AlterSql(model Table, diff TableDiff) ([]byte, error)
	DeleteSql(model Table) ([]byte, error)
	ParseCreateQuery(model *Table, query string) error
	ParseUpdateQuery(model *Table, query string) error
}
