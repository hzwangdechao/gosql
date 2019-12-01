package connect

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	// "math"
	"reflect"
	"strings"
	"time"
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
	return &Connect{db, db != nil && db.Ping() != nil}
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
			transferType(col, types[idx], value)
		}
		result = append(result, col)
	}
	rows.Close()
	return
}

func transferType(m Result, t *sql.ColumnType, value any) {
	name := t.Name()
	switch t.DatabaseTypeName() {
	case "VARCHAR", "TEXT":
		m[name] = string(value.([]byte))
	case "INT":
		n, _ := strconv.ParseInt(string(value.([]byte)), 10, 32)
		m[name] = int32(n)
	case "BIGINT":
		n, _ := strconv.ParseInt(string(value.([]byte)), 10, 64)
		m[name] = int64(n)
	case "SMALLINT":
		n, _ := strconv.ParseInt(string(value.([]byte)), 10, 16)
		m[name] = int16(n)
	case "TINYINT":
		n, _ := strconv.ParseInt(string(value.([]byte)), 10, 8)
		m[name] = int8(n)
	case "BIT":
		m[name] = value.([]byte)[0]
	case "DATE":
		m[name], _ = time.Parse("2006-01-02", string(value.([]byte)))
	case "DATETIME":
		m[name], _ = time.Parse("2006-01-02 15:04:05", string(value.([]byte)))
	case "DOUBLE", "DECIMAL":
		m[name], _ = strconv.ParseFloat(string(value.([]byte)), 64)
	case "FLOAT":
		f64, _ := strconv.ParseFloat(string(value.([]byte)), 32)
		m[name] = float32(f64)
	default:
		fmt.Println(t.DatabaseTypeName(), reflect.TypeOf(value).Name(), value)
		m[name] = value
	}
}
