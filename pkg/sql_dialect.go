package pkg

import (
	"fmt"
	"strings"
)

type sqlDialect interface {
	CreateSql(model apiModel) (string, error)
	AddFieldsSql(model apiModel, fields []apiModelField) (string, error)
	DeleteFieldsSql(model apiModel, fields []apiModelField) (string, error)
	UpdateFieldsSql(model apiModel, fields []apiModelField) (string, error)
	DeleteSql(model apiModel) (string, error)
	ParseCurrentModels(path string) ([]apiModel, error)
}

var _ sqlDialect = &mysqlDialect{}

type mysqlDialect struct {
}

func newMysqlDialect() *mysqlDialect {
	return &mysqlDialect{}
}

func (d *mysqlDialect) CreateSql(model apiModel) (string, error) {
	createTableWrapper := d.getModelCreateTableWrapper(model)

	columnsSql := ""
	for _, field := range model.Fields {
		columnSql, err := d.getColumnSql(field)
		if err != nil {
			return "", fmt.Errorf("failed to generate sql to create column %s.%s :%w", model.Name, field.Name, err)
		}

		columnsSql += fmt.Sprintf("%s,\n", columnSql)
	}

	indexesSql := ""
	separator := ",\n"
	for _, field := range model.Fields {
		for _, c := range field.Constraints {
			indexSql := ""
			switch c.Name {
			case "foreign_key":
				keyName := getFieldSqlColumn(field.Name)
				forTable := getModelSqlTable(c.Value1)
				forColumn := getFieldSqlColumn(c.Value2)
				indexSql = fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)", keyName, keyName, forTable, forColumn)

			default:
				// do nothing
				continue
			}

			indexesSql += fmt.Sprintf("%s\n%s", separator, indexSql)
			separator = ","
		}
	}

	return fmt.Sprintf(createTableWrapper, columnsSql, indexesSql), nil
}

func (d *mysqlDialect) AddFieldsSql(model apiModel, fields []apiModelField) (string, error) {
	updateTableWrapper := d.getModelUpdateTableWrapper(model)

	columnsSql := ""
	separator := ""
	for _, field := range fields {
		columnSql, err := d.getColumnSql(field)
		if err != nil {
			return "", fmt.Errorf("failed to generate sql to add column %s.%s :%w", model.Name, field.Name, err)
		}

		columnsSql += fmt.Sprintf("%sADD COLUMN %s", separator, columnSql)
		separator = ",\n"
	}

	separator = ",\n\n"
	for _, field := range model.Fields {
		for _, c := range field.Constraints {
			indexSql := ""
			switch c.Name {
			case "foreign_key":
				keyName := getFieldSqlColumn(field.Name)
				forTable := getModelSqlTable(c.Value1)
				forColumn := getFieldSqlColumn(c.Value2)
				indexSql = fmt.Sprintf("ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)", keyName, keyName, forTable, forColumn)

			default:
				// do nothing
				continue
			}

			columnsSql += fmt.Sprintf("%s%s", separator, indexSql)
			separator = ",\n"
		}
	}

	return fmt.Sprintf(updateTableWrapper, columnsSql), nil
}

func (d *mysqlDialect) UpdateFieldsSql(model apiModel, fields []apiModelField) (string, error) {
	updateTableWrapper := d.getModelUpdateTableWrapper(model)

	columnsSql := ""
	separator := ""
	for _, field := range fields {
		columnSql, err := d.getColumnSql(field)
		if err != nil {
			return "", fmt.Errorf("failed to generate sql to add column %s.%s :%w", model.Name, field.Name, err)
		}

		columnsSql += fmt.Sprintf("%sMODIFY COLUMN %s", separator, columnSql)
		separator = ",\n"
	}

	return fmt.Sprintf(updateTableWrapper, columnsSql), nil
}

func (d *mysqlDialect) DeleteFieldsSql(model apiModel, fields []apiModelField) (string, error) {
	updateTableWrapper := d.getModelUpdateTableWrapper(model)

	columnsSql := ""
	separator := ""
	for _, field := range fields {
		columnsSql += fmt.Sprintf("%sDROP COLUMN %s", separator, getFieldSqlColumn(field.Name))
		separator = ",\n"
	}

	separator = ",\n\n"
	for _, field := range model.Fields {
		indexSql := ""
		for _, c := range field.Constraints {
			switch c.Name {
			case "foreign_key":
				indexSql = fmt.Sprintf("DROP FOREIGN KEY %s", getFieldSqlColumn(field.Name))

			default:
				// do nothing
				continue
			}
			columnsSql += fmt.Sprintf("%s%s", separator, indexSql)
			separator = ",\n"
		}
	}

	return fmt.Sprintf(updateTableWrapper, columnsSql), nil
}

func (d *mysqlDialect) DeleteSql(model apiModel) (string, error) {
	return fmt.Sprintf("drop table %s;", getModelSqlTable(model.Name)), nil
}

func (d *mysqlDialect) ParseCurrentModels(path string) ([]apiModel, error) {
	panic("implement me")
}

func (d *mysqlDialect) getModelCreateTableWrapper(model apiModel) string {
	return fmt.Sprintf(`CREATE TABLE %s
(
id BIGINT AUTO_INCREMENT PRIMARY KEY,
uuid VARCHAR(36) NOT NULL UNIQUE,

%%s
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP%%s
);`, getModelSqlTable(model.Name))
}

func (d *mysqlDialect) getModelUpdateTableWrapper(model apiModel) string {
	return fmt.Sprintf(`ALTER TABLE %s
%%s;`, getModelSqlTable(model.Name))
}

func (d *mysqlDialect) getColumnSql(field apiModelField) (string, error) {
	typeString := field.Type

	// if type:example constraint set type to provided value
	if c, ok := field.getConstraint("type"); ok {
		typeString = c.Value1
	}

	sqlColumn := getFieldSqlColumn(field.Name)

	switch strings.ReplaceAll(typeString, "*", "") {
	case "string":
		sqlColumn += " VARCHAR(255)"
	case "bool":
		sqlColumn += " BOOL"
	case "int64":
		sqlColumn += " BIGINT"
	case "int":
		sqlColumn += " INT"
	case "text":
		sqlColumn += " TEXT"
	default:
		return "", fmt.Errorf("unknown type %s", typeString)
	}

	if !strings.HasPrefix(typeString, "*") {
		sqlColumn += " NOT NULL"
	}

	if _, ok := field.getConstraint("unique"); ok {
		sqlColumn += " UNIQUE"
	}

	if c, ok := field.getConstraint("default"); ok {
		sqlColumn += fmt.Sprintf(" DEFAULT %s", c.Value1)
	}

	return sqlColumn, nil
}
