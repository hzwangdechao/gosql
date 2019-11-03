package sql

import (
	"fmt"
)

type Table struct {
	Database string
	Name     string
	Alias    string
}

func NewTable(database string, name string, alias string) *Table {
	return &Table{
		Database: database,
		Name:     name,
		Alias:    alias,
	}
}

func (tb *Table) SQL() *SQL {
	return NewSQL().From(tb)
}

func (tb *Table) As(alias string) *Table {
	return NewTable(tb.Database, tb.Name, alias)
}

// String return from part of table
func (tb *Table) String() string {
	if tb.Alias != "" {
		return fmt.Sprintf("`%s`.`%s` AS `%s`", tb.Database, tb.Name, tb.Alias)
	}
	return fmt.Sprintf("`%s`.`%s`", tb.Database, tb.Name)
}

func (tb *Table) Prefix() string {
	if tb.Alias != "" {
		return fmt.Sprintf("`%s`", tb.Alias)
	}
	return fmt.Sprintf("`%s`.`%s`", tb.Database, tb.Name)
}
