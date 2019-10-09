package sql

import (
	"fmt"
)

// Column of database
type Column struct {
	Table *Table
	Name  string
	Alias string
	Type  Type
}

func NewColumn(table *Table, name string, alias string, columnType Type) *Column {
	return &Column{
		Table: table,
		Name:  name,
		Alias: alias,
		Type:  columnType,
	}
}

// String return the Column as string
func (col *Column) String() string {
	if col.Type == ColumnTypeFunction {
		if col.Alias != "" {
			return fmt.Sprintf("`%s` as `%s`", col.Name, col.Alias)
		}
		return fmt.Sprintf("`%s`", col.Name)
	}
	if col.Alias != "" {
		return fmt.Sprintf("%s.`%s` as `%s`", col.Table.String(), col.Name, col.Alias)
	}
	return fmt.Sprintf("%s.`%s`", col.Table.String(), col.Name)
}
