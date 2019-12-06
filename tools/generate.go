// TestFlag project main.go
package main

import (
	"flag"
	"fmt"
	"github.com/OhYee/gosql/utils/connect"
	"go/format"
	"os"
	"strings"
)

const (
	template = `
package %s

import (
	"github.com/OhYee/gosql"
)

type %s struct {
	*sql.Table
	%s
}

func New%s(alias ...string) *%s {
	const size = %d
	as := make([]string, size)
	for i := 0; i < size; i++ {
		if len(alias) >= i+1 {
			as[i] = alias[i]
		} else {
			as[i] = ""
		}
	}

	tb := sql.NewTable("%s", "%s", as[0])
	return &%s{
		Table:    tb,
		%s
	}
}
`
)

func main() {
	var port int
	var output, pkg, ip, user, password, database, table string

	flag.StringVar(&output, "output", "", "output file")
	flag.StringVar(&output, "o", "", "output file")

	flag.StringVar(&pkg, "package", "sql", "output package")
	flag.StringVar(&pkg, "pkg", "sql", "output package")

	flag.StringVar(&user, "user", "root", "database user")
	flag.StringVar(&user, "u", "root", "database user")

	flag.StringVar(&password, "password", "", "database password")
	flag.StringVar(&password, "pw", "", "database password")

	flag.StringVar(&ip, "ip", "127.0.0.1", "database port")
	flag.StringVar(&ip, "i", "127.0.0.1", "database port")

	flag.IntVar(&port, "port", 3306, "database port")
	flag.IntVar(&port, "p", 3306, "database port")

	flag.StringVar(&database, "database", "", "database name")
	flag.StringVar(&database, "d", "", "database name")

	flag.StringVar(&table, "table", "", "table name")
	flag.StringVar(&table, "t", "", "table name")

	flag.Parse()

	conn, err := connect.NewConnection(&connect.DataSourceName{Username: user, Password: password, Address: ip, Port: port})
	if err != nil {
		panic(err)
	}
	res, err := conn.Query("select `COLUMN_NAME`, `DATA_TYPE` from information_schema.columns where table_schema=? and table_name=? ORDER BY `ORDINAL_POSITION`", database, table)
	if err != nil {
		panic(err)
	}

	type column struct {
		name string
		t    string
	}
	columns := make([]column, len(res))

	for i, r := range res {
		var name, t interface{}
		var ok bool
		if name, ok = r["COLUMN_NAME"]; !ok {
			continue
		}
		if t, ok = r["COLUMN_NAME"]; !ok {
			continue
		}
		columns[i] = column{
			name: name.(string),
			t:    t.(string),
		}
	}

	upper := func(s string) string {
		if len(s) == 0 {
			return s
		}
		return strings.ToUpper(s[0:1]) + s[1:]
	}

	tableUpper := upper(table)
	columnStruct := make([]string, len(columns))
	columnDeclare := make([]string, len(columns))
	for i, col := range columns {
		columnStruct[i] = fmt.Sprintf("%s\t*sql.Column", upper(col.name))
		columnDeclare[i] = fmt.Sprintf("%s: sql.NewColumn(tb, \"%s\", as[%d], sql.ColumnTypeUnknown),", upper(col.name), col.name, i+1)
	}

	var source []byte
	temp := []byte(fmt.Sprintf(template, pkg, tableUpper, strings.Join(columnStruct, "\n\t"), tableUpper, tableUpper, len(columns), database, table, tableUpper, strings.Join(columnDeclare, "\n\t\t")))
	if source, err = format.Source(temp); err != nil {
		source = temp
	}

	if output != "" {
		outputFile, err := os.Create(output)
		if err == nil {
			outputFile.Write(source)
			outputFile.Close()
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("%s\n", source)
	}

}
