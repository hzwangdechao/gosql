package sql

import (
	"fmt"
	"strings"
)

// SQL
type SQL struct {
	columns    []*Column
	tables     []*Table
	conditions []interface{}
}

func NewSQL() *SQL {
	return &SQL{
		columns:    make([]*Column, 0),
		conditions: make([]interface{}, 0),
		tables:     make([]*Table, 0),
	}
}

func (sql *SQL) Select(columns ...*Column) *SQL {
	for _, column := range columns {
		sql.columns = append(sql.columns, column)
	}
	return sql
}

func (sql *SQL) From(tables ...*Table) *SQL {
	for _, table := range tables {
		sql.tables = append(sql.tables, table)
	}
	return sql
}

func (sql *SQL) Query() string {
	return strings.Join([]string{
		sql.getSelectPart(),
		sql.getFromPart(),
	}, " ") + ";"
}

func (sql *SQL) getSelectPart() string {
	columns := make([]string, 0)
	for _, column := range sql.columns {
		columns = append(columns, column.String())
	}
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	return fmt.Sprintf("SELECT %s", strings.Join(columns, ", "))
}

func (sql *SQL) getFromPart() string {
	tables := make([]string, 0)
	for _, table := range sql.tables {
		tables = append(tables, table.String())
	}
	return fmt.Sprintf("FROM %s", strings.Join(tables, " and "))
}

