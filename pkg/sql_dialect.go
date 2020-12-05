package pkg

import (
	"fmt"
	"strings"
)

type sqlDialect interface {
	CreateSql(model apiModel) (string, error)
	AddFieldsSql(model apiModel, newFields []apiModelField) (string, error)
	DeleteFieldsSql(model apiModel, deletedFields []apiModelField) (string, error)
	DeleteSql(model apiModel) (string, error)
	ParseCreateQuery(model *apiModel, query string) error
	ParseUpdateQuery(model *apiModel, query string) error
}

var _ sqlDialect = &mysqlDialect{}

type mysqlDialect struct {
}

func newMysqlDialect() *mysqlDialect {
	return &mysqlDialect{}
}

var mysqlTypeMappings = map[string]string{
	"string":    "VARCHAR(255)",
	"bool":      "BOOL",
	"int64":     "BIGINT",
	"int":       "INT",
	"text":      "TEXT",
	"timestamp": "TIMESTAMP",
}

func (d *mysqlDialect) CreateSql(model apiModel) (string, error) {
	var statements []string
	for _, field := range model.Fields {
		statement, err := d.getColumnSql(field)
		if err != nil {
			return "", fmt.Errorf("failed to generate column sql %s.%s :%w", model.Name, field.Name, err)
		}
		statements = append(statements, statement)
	}

	for _, field := range model.Fields {
		for _, c := range field.Constraints {
			statement, ok := d.getConstraintSql(field, c)
			if ok {
				statements = append(statements, statement)
			}
		}
	}

	return fmt.Sprintf(`CREATE TABLE %s
(
%s
);`, getModelSqlTable(model.Name), strings.Join(statements, ",\n")), nil
}

func (d *mysqlDialect) AddFieldsSql(model apiModel, newFields []apiModelField) (string, error) {
	var statements []string
	for _, oldField := range newFields {
		field := model.mustGetField(oldField.Name)
		statement, err := d.getColumnSql(field)
		if err != nil {
			return "", fmt.Errorf("failed to generate column sql %s.%s :%w", model.Name, field.Name, err)
		}
		statements = append(statements, fmt.Sprintf("ADD COLUMN %s", statement))
	}

	for _, oldField := range newFields {
		field := model.mustGetField(oldField.Name)
		for _, c := range field.Constraints {
			statement, ok := d.getConstraintSql(field, c)
			if ok {
				statements = append(statements, fmt.Sprintf("ADD %s", statement))
			}
		}
	}

	return fmt.Sprintf(`ALTER TABLE %s
%s;`, getModelSqlTable(model.Name), strings.Join(statements, ",\n")), nil
}

func (d *mysqlDialect) DeleteFieldsSql(model apiModel, deletedFields []apiModelField) (string, error) {
	var statements []string
	for _, field := range deletedFields {
		statement := fmt.Sprintf("DROP COLUMN %s", getFieldSqlColumn(field.Name))
		statements = append(statements, statement)
	}

	for _, field := range deletedFields {
		for _, c := range field.Constraints {
			switch c.Name {
			case "foreign_key":
				statement := fmt.Sprintf("DROP FOREIGN KEY %s", getFieldSqlColumn(field.Name))
				statements = append(statements, statement)
			}
		}
	}

	return fmt.Sprintf(`ALTER TABLE %s
%s;`, getModelSqlTable(model.Name), strings.Join(statements, ",\n")), nil
}

func (d *mysqlDialect) DeleteSql(model apiModel) (string, error) {
	return fmt.Sprintf("drop table %s;", getModelSqlTable(model.Name)), nil
}

func (d *mysqlDialect) getColumnSql(field apiModelField) (string, error) {
	typeString := field.Type

	// if type:example constraint set type to provided value
	if c, ok := field.GetConstraint("type"); ok {
		typeString = c.Value1
	}

	columnName := getFieldSqlColumn(field.Name)
	sqlType, err := d.convertTypeToSqlType(strings.ReplaceAll(typeString, "*", ""))
	if err != nil {
		return "", err
	}

	var modifiers []string
	if !strings.HasPrefix(typeString, "*") {
		modifiers = append(modifiers, "NOT NULL")
	}

	if field.HasConstraint("unique") {
		modifiers = append(modifiers, "UNIQUE")
	}

	if field.HasConstraint("auto_increment") {
		modifiers = append(modifiers, "AUTO_INCREMENT")
	}

	if field.HasConstraint("primary_key") {
		modifiers = append(modifiers, "PRIMARY KEY")
	}

	if c, ok := field.GetConstraint("default"); ok {
		modifiers = append(modifiers, fmt.Sprintf("DEFAULT %s", c.Value1))
	}

	if c, ok := field.GetConstraint("on_update"); ok {
		modifiers = append(modifiers, fmt.Sprintf("ON UPDATE %s", c.Value1))
	}

	combinedModifiers := strings.TrimSpace(strings.Join(modifiers, " "))
	return strings.TrimSpace(fmt.Sprintf("%s %s %s", columnName, sqlType, combinedModifiers)), nil
}

func (d *mysqlDialect) getConstraintSql(field apiModelField, c constraint) (string, bool) {
	switch c.Name {
	case "foreign_key":
		keyName := getFieldSqlColumn(field.Name)
		forTable := getModelSqlTable(c.Value1)
		forColumn := getFieldSqlColumn(c.Value2)
		return fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)", keyName, keyName, forTable, forColumn), true
	}
	return "", false
}

func (d *mysqlDialect) convertTypeToSqlType(fieldType string) (string, error) {
	sqlType, ok := mysqlTypeMappings[fieldType]
	if !ok {
		return "", fmt.Errorf("could not find matching sql type for %s", fieldType)
	}
	return sqlType, nil
}

func (d *mysqlDialect) convertSqlTypeToFieldType(sqlType string) (string, error) {
	for k, v := range mysqlTypeMappings {
		if v == sqlType {
			return k, nil
		}
	}
	return "", fmt.Errorf("could not find matching type for sql %s", sqlType)
}

func (d *mysqlDialect) ParseCreateQuery(model *apiModel, query string) error {
	query = strings.ReplaceAll(strings.ReplaceAll(query, ";", ""), ",", "")
	lines := strings.Split(query, "\n")
	lines = lines[2 : len(lines)-1]

	fieldsMap := map[string]apiModelField{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		words := strings.Split(line, " ")
		if len(words) < 2 {
			return fmt.Errorf("%s does not seem like a valid field sql", line)
		}

		// can be a field definition or a constraint definition
		switch words[0] {
		case "CONSTRAINT":
			fieldName := getFieldNameFromSqlColumn(words[1])
			field, ok := fieldsMap[fieldName]
			if !ok {
				return fmt.Errorf("found constraint for missing field %s", fieldName)
			}
			err := d.parseFieldSqlConstraint(&field, line)
			if err != nil {
				return fmt.Errorf("failed to parse constraint %s for field %s: %w", line, field.Name, err)
			}
			fieldsMap[field.Name] = field
		default:
			// must be a field definition
			field := apiModelField{}
			err := d.parseFieldSql(&field, line)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", line, err)
			}
			fieldsMap[field.Name] = field
		}
	}

	for _, field := range fieldsMap {
		model.Fields = append(model.Fields, field)
	}
	return nil
}

func (d *mysqlDialect) ParseUpdateQuery(model *apiModel, query string) error {
	query = strings.ReplaceAll(strings.ReplaceAll(query, ";", ""), ",", "")
	lines := strings.Split(query, "\n")
	lines = lines[2:]

	fieldsMap := map[string]apiModelField{}
	for _, field := range model.Fields {
		fieldsMap[field.Name] = field
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		words := strings.Split(line, " ")
		if len(words) < 3 {
			return fmt.Errorf("%s does not seem like a valid update field sql", line)
		}

		firstWord := words[0]
		secondWord := words[1]
		thirdWord := words[2]
		forthWord := ""
		if len(words) > 3 {
			forthWord = words[3]
		}
		firstTwoWords := fmt.Sprintf("%s %s", firstWord, secondWord)

		// can be a field definition or a constraint definition
		switch firstTwoWords {
		case "ADD COLUMN":
			columnSql := strings.ReplaceAll(line, "ADD COLUMN ", "")
			field := apiModelField{}
			err := d.parseFieldSql(&field, columnSql)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", line, err)
			}
			fieldsMap[field.Name] = field
		case "ADD CONSTRAINT":
			fieldName := getFieldNameFromSqlColumn(thirdWord)
			field, ok := fieldsMap[fieldName]
			if !ok {
				return fmt.Errorf("found constraint for missing field %s", fieldName)
			}
			constraintSql := strings.ReplaceAll(line, "ADD ", "")
			err := d.parseFieldSqlConstraint(&field, constraintSql)
			if err != nil {
				return fmt.Errorf("failed to parse constraint %s for field %s: %w", line, field.Name, err)
			}
			fieldsMap[field.Name] = field
		case "DROP COLUMN":
			delete(fieldsMap, getFieldNameFromSqlColumn(thirdWord))
		case "DROP FOREIGN KEY":
			fieldName := getFieldNameFromSqlColumn(forthWord)
			field, ok := fieldsMap[fieldName]
			if !ok {
				// just skip the constraint
				continue
			}
			var constraints []constraint
			for _, c := range field.Constraints {
				if c.Name == "foreign_key" {
					continue
				}
				constraints = append(constraints, c)
			}
			field.Constraints = constraints
		}
	}

	var fields []apiModelField
	for _, field := range model.Fields {
		fields = append(fields, field)
	}
	return nil
}

func (d *mysqlDialect) parseFieldSql(field *apiModelField, fieldSql string) error {
	words := strings.Split(fieldSql, " ")
	if len(words) < 2 {
		return fmt.Errorf("%s does not seem like a valid field sql", fieldSql)
	}

	field.Name = getFieldNameFromSqlColumn(words[0])
	fieldType, err := d.convertSqlTypeToFieldType(words[1])
	if err != nil {
		return fmt.Errorf("failed to get field type: %w", err)
	}
	field.Type = fieldType

	words = words[2:]
	for i := 0; i < len(words); i++ {
		firstWord := words[i]
		secondWord := ""
		if len(words) > i+1 {
			secondWord = words[i+1]
		}
		firstTwoWords := fmt.Sprintf("%s %s", firstWord, secondWord)
		thirdWord := ""
		if len(words) > i+2 {
			thirdWord = words[i+2]
		}

		switch firstTwoWords {
		case "NOT NULL":
			fieldType = "*" + fieldType
			continue
		case "ON UPDATE":
			field.Constraints = append(field.Constraints, constraint{Name: "on_update", Value1: thirdWord})
			continue
		case "PRIMARY KEY":
			field.Constraints = append(field.Constraints, constraint{Name: "primary_key"})
			continue
		}

		switch firstWord {
		case "DEFAULT":
			field.Constraints = append(field.Constraints, constraint{Name: "default", Value1: secondWord})
			continue
		case "AUTO_INCREMENT":
			field.Constraints = append(field.Constraints, constraint{Name: "auto_increment"})
			continue
		case "UNIQUE":
			field.Constraints = append(field.Constraints, constraint{Name: "unique"})
			continue
		}
	}
	return nil
}

func (d *mysqlDialect) parseFieldSqlConstraint(field *apiModelField, constraintSql string) error {
	// here it does not matter if its an alter table or create table or column creation or modification
	constraintSql = strings.ReplaceAll(constraintSql, "CONSTRAINT", "")
	constraintSql = strings.TrimSpace(constraintSql)

	words := strings.Split(constraintSql, " ")[1:] //dont need first word its sql field name
	firstWord := words[0]
	secondWord := ""
	if len(words) > 1 {
		secondWord = words[1]
	}
	firstTwoWords := fmt.Sprintf("%s %s", firstWord, secondWord)

	switch firstTwoWords {
	case "FOREIGN KEY":
		foreignTableSqlName := words[len(words)-2]
		foreignFieldSqlName := strings.Trim(words[len(words)-1], " ()")
		field.Constraints = append(field.Constraints, constraint{
			Name:   "foreign_key",
			Value1: getModelNameFromSqlTable(foreignTableSqlName),
			Value2: getFieldNameFromSqlColumn(foreignFieldSqlName),
		})
	}

	return nil
}
