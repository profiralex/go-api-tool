package migrations

import (
	"bytes"
	"github.com/profiralex/go-api-tool/pkg/files"
	"testing"
)

func TestMysqlDialectImplementsSqlDialect(t *testing.T) {
	var _ SqlDialect = &MysqlDialect{}
}

func TestMysqlDialectCreate(t *testing.T) {
	var d SqlDialect = &MysqlDialect{fp: files.NewBoxProvider()}
	model := Table{
		Name: "records",
		Columns: []Column{
			{
				Name:          "id",
				PrimaryKey:    true,
				AutoIncrement: true,
				Type:          "BIGINT",
			},
			{
				Name:          "user_id",
				PrimaryKey:    false,
				AutoIncrement: false,
				Nullable:      true,
				Type:          "BIGINT",
			},
			{
				Name:          "email",
				PrimaryKey:    false,
				AutoIncrement: false,
				Nullable:      true,
				Type:          "VARCHAR(255)",
			},
			{
				Name:          "updated_at",
				PrimaryKey:    false,
				AutoIncrement: false,
				Nullable:      true,
				Type:          "TIMESTAMP",
				DefaultExpr:   "CURRENT_TIMESTAMP",
				OnUpdateExpr:  "CURRENT_TIMESTAMP",
			},
		},
		ForeignKeys: []ForeignKey{
			{
				Field:        "user_id",
				TableName:    "users",
				ForeignField: "id",
			},
		},
		Indexes: []Index{
			{
				Name:    "unique",
				Columns: []string{"user_id", "email"},
			},
		},
	}
	expectedSql := []byte(`CREATE TABLE records (
CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES users (id),
UNIQUE KEY unique (user_id, email),
id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
user_id BIGINT,
email VARCHAR(255),
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);`)

	sql, err := d.CreateSql(model)
	if err != nil {
		t.Errorf("failed to create sql: %s", err)
		return
	}

	if bytes.Compare(expectedSql, sql) != 0 {
		t.Errorf("assert failed\nExpected:\n%s\nGot:\n%s", expectedSql, sql)
		return
	}
}
