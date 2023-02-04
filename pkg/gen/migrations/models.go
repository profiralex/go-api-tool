package migrations

import (
	"fmt"
)

type Table struct {
	Name        string
	Columns     []Column
	ForeignKeys []ForeignKey
	Indexes     []Index
}

type Column struct {
	Name          string
	PrimaryKey    bool
	AutoIncrement bool
	Nullable      bool
	Type          string
	DefaultExpr   string
	OnUpdateExpr  string
}

type ForeignKey struct {
	Field        string
	TableName    string
	ForeignField string
}

type Index struct {
	Name    string
	Columns []string
}

type TableDiff struct {
	AddedColumns   []Column
	UpdatedColumns []Column
	DroppedColumns []Column

	AddedForeignKeys   []ForeignKey
	DroppedForeignKeys []ForeignKey

	AddedIndexes   []Index
	DroppedIndexes []Index
}

func (d *TableDiff) HasAddedColumns() bool {
	return len(d.AddedColumns) > 0
}

func (d *TableDiff) HasUpdatedColumns() bool {
	return len(d.UpdatedColumns) > 0
}

func (d *TableDiff) HasDroppedColumns() bool {
	return len(d.DroppedColumns) > 0
}

func (d *TableDiff) HasAddedForeignKeys() bool {
	return len(d.AddedForeignKeys) > 0
}

func (d *TableDiff) HasDroppedForeignKeys() bool {
	return len(d.DroppedForeignKeys) > 0
}

func (d *TableDiff) HasAddedIndexes() bool {
	return len(d.AddedIndexes) > 0
}

func (d *TableDiff) HasDroppedIndexes() bool {
	return len(d.DroppedIndexes) > 0
}

func (d *TableDiff) HasAddedData() bool {
	return d.HasAddedColumns() || d.HasUpdatedColumns() || d.HasAddedForeignKeys() || d.HasAddedIndexes()
}

func (d *TableDiff) HasDroppedData() bool {
	return d.HasDroppedColumns() || d.HasDroppedForeignKeys() || d.HasDroppedIndexes()
}

type apiModelUpdate struct {
	Table               Table
	ColumnsToAdd        []Column
	ColumnsToUpdate     []Column
	ColumnsToDelete     []Column
	ForeignKeysToAdd    []ForeignKey
	ForeignKeysToUpdate []ForeignKey
	ForeignKeysToDelete []ForeignKey
}

func (m *Table) GetField(columnName string) (Column, bool) {
	for _, field := range m.Columns {
		if field.Name == columnName {
			return field, true
		}
	}

	return Column{}, false
}

func (m *Table) MustGetField(columnName string) Column {
	field, ok := m.GetField(columnName)
	if !ok {
		panic(fmt.Sprintf("failed to get field %s from table %s", columnName, m.Name))
	}
	return field
}

func (m *Table) HasField(columnName string) bool {
	_, ok := m.GetField(columnName)
	return ok
}
