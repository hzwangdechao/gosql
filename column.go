package sql

import (
	"fmt"
	"github.com/OhYee/goutils/transfer"
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

func NewFunctionColumn(function string, alias string, args ...string) *Column {
	return &Column{
		Table: nil,
		Name:  fmt.Sprintf(function, transfer.ToInterfaceSlice(args)...),
		Alias: alias,
		Type:  ColumnTypeFunction,
	}
}

func (col *Column) As(alias string) *Column {
	return NewColumn(col.Table, col.Name, alias, col.Type)
}

// String return the Column as string
func (col *Column) String() string {
	if col.Type == ColumnTypeFunction {
		if col.Alias != "" {
			return fmt.Sprintf("%s AS `%s`", col.Name, col.Alias)
		}
		return fmt.Sprintf("%s", col.Name)
	}
	if col.Alias != "" {
		return fmt.Sprintf("%s.`%s` AS `%s`", col.Table.Prefix(), col.Name, col.Alias)
	}
	return fmt.Sprintf("%s.`%s`", col.Table.Prefix(), col.Name)
}
