package sql

import (
	"fmt"
	"github.com/OhYee/gosql/operator"
	"github.com/OhYee/goutils/functional"
	"strings"
)

type any = interface{}

// SQL
type SQL struct {
	columns    []*Column
	tables     []*Table
	conditions []*op.Operator
}

func NewSQL() *SQL {
	return &SQL{
		columns:    make([]*Column, 0),
		tables:     make([]*Table, 0),
		conditions: make([]*op.Operator, 0),
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

func (sql *SQL) Where(conditions ...*op.Operator) *SQL {
	for _, op := range conditions {
		sql.conditions = append(sql.conditions, op)
	}
	return sql
}

// Query return the string of the sql query (for send to server, will add semicolon)
func (sql *SQL) Query() string {
	return sql.toString() + ";"
}

// toString return the string of the sql query (without brackets and semicolon)
func (sql *SQL) toString() string {
	strSlice := []string{
		sql.getSelectPart(),
		sql.getFromPart(),
		sql.getWherePart(),
	}

	strSlice = fp.FilterString(func(s string) bool {
		return len(s) > 0
	}, strSlice)

	return strings.Join(strSlice, " ")
}

// ToString return the string of this sql query (for sub-query, will add brackets)
func (sql *SQL) ToString() string {
	return fmt.Sprintf("(%s)", sql.toString())
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
	return fmt.Sprintf("FROM %s", strings.Join(tables, ", "))
}

func (sql *SQL) getWherePart() string {
	if len(sql.conditions) == 0 {
		return ""
	}
	conditions := make([]string, 0)
	for _, condition := range sql.conditions {
		conditions = append(conditions, condition.String())
	}
	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " and "))
}

func (sql *SQL) Add(equations ...string) string {
	tables := make([]string, 0)
	for _, table := range sql.tables {
		tables = append(tables, table.String())
	}

	columns := make([]string, 0)
	values := make([]string, 0)

	for _, equation := range equations {
		column := strings.Join(strings.Split(equation, "=")[:1], "")
		value := strings.Join(strings.Split(equation, "=")[1:], "=")
		fmt.Println(column)
		columns = append(columns, column)
		values = append(values, value)
	}
	if len(columns) > 0 {
		strColumn := fmt.Sprintf("`%s`", strings.Join(columns, "` ,`"))
		strValues := fmt.Sprintf("'%s'", strings.Join(values, "' ,'"))
		return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", strings.Join(tables, " "), strColumn, strValues)

	} else {
		panic("Expression cannot be empty")
	}

}

func (sql *SQL) Delete(configure ...bool) string {
	// configure 防止因为没有 where 条件导致删除整个数据库
	// 在没有where条件时需要传入一个true确认是自己想要执行的命令
	if sql.getWherePart() == "" && len(configure) < 1 {
		panic("Operation not allowed")
	} else {
		return fmt.Sprintf("DELETE  %s %s", sql.getFromPart(), sql.getWherePart())
	}
}

func (sql *SQL) Update(equations ...string) string {
	tables := make([]string, 0)
	for _, table := range sql.tables {
		tables = append(tables, table.String())
	}

	if sql.getWherePart() != "" {
		return fmt.Sprintf("UPDATE %s SET %s %s", strings.Join(tables, " "), strings.Join(equations, ","), sql.getWherePart())

	} else {
		panic("Expression cannot be empty")
	}

}
