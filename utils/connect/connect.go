package connect

import (
	"database/sql"
	"fmt"
	// "github.com/go-sql-driver/mysql"
	"strings"
)

type any = interface{}
type Result = map[string]any
type Results = []Result

// DataSourceName database source name
type DataSourceName struct {
	Username string
	Password string
	Address  string
	Port     int
}

func (dsn *DataSourceName) String() string {
	s := make([]string, 0)
	if dsn.Username != "" {
		s = append(s, dsn.Username)
	}
	if dsn.Password != "" {
		s = append(s, ":"+dsn.Password)
	}
	if dsn.Address != "" {
		ad := dsn.Address
		if dsn.Port != 0 {
			ad = fmt.Sprintf("%s:%d", ad, dsn.Port)
		}
		s = append(s, "@tcp("+ad+")")
	}
	return strings.Join(s, "") + "/"
}

// Connect a database connection
type Connect struct {
	db      *sql.DB
	isClose bool
}

func newConnect(db *sql.DB) *Connect {
	return &Connect{db, db.Ping() != nil}
}

// NewConnection connect to a exist database
func NewConnection(dsn *DataSourceName) (conn *Connect, err error) {
	db, err := sql.Open("mysql", dsn.String())
	conn = newConnect(db)
	return
}

// Close the database connects
func (conn *Connect) Close() (err error) {
	if !conn.isClose {
		err = conn.db.Close()
		if err == nil {
			conn.isClose = true
		}
	}
	return
}

// Query the data according the sql string
func (conn *Connect) Query(sqlStr string, args ...any) (result Results, err error) {
	if conn.isClose {
		err = fmt.Errorf("Database connect is closed")
		return
	}
	rows, err := conn.db.Query(sqlStr, args...)
	if err != nil {
		return
	}

	result = make(Results, 0)
	types, err := rows.ColumnTypes()
	if err != nil {
		return
	}

	for rows.Next() {
		col := make(Result)
		tmp := make([]any, len(types))
		for idx := range tmp {
			tmp[idx] = &tmp[idx]
		}

		err = rows.Scan(tmp...)
		if err != nil {
			return
		}
		for idx, value := range tmp {
			col[types[idx].Name()] = value
		}
		result = append(result, col)
	}
	rows.Close()
	return
}
