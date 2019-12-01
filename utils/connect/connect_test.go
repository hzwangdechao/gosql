package connect

import (
	"github.com/OhYee/gosql"
	"github.com/OhYee/goutils"
	"os"
	"testing"
	"time"
)

var validDatabase = DataSourceName{"ohyee", "123456", "127.0.0.1", 3306}
var isLocal = func() bool {
	hostname, _ := os.Hostname()
	return hostname == "OhYee-wsl"
}()

type User struct {
	// *SQL
	*sql.Table
	Username *sql.Column
	Realname *sql.Column
}

func NewUser(alias ...string) *User {
	const size = 3
	as := make([]string, size)
	for i := 0; i < size; i++ {
		if len(alias) >= i+1 {
			as[i] = alias[i]
		} else {
			as[i] = ""
		}
	}

	tb := sql.NewTable("test", "users", as[0])
	return &User{
		// SQL:      NewSQL().From(tb),
		Table:    tb,
		Username: sql.NewColumn(tb, "username", as[1], sql.ColumnTypeUnknown),
		Realname: sql.NewColumn(tb, "realname", as[2], sql.ColumnTypeUnknown),
	}
}

func TestNewConnectionAndClose(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name           string
		dsn            *DataSourceName
		isClose        bool
		wantErr        bool
		wantCloseError bool
	}{
		{
			name:           "Connect to a exist database",
			dsn:            &validDatabase,
			isClose:        !isLocal,
			wantErr:        false,
			wantCloseError: !isLocal,
		},
		{
			name:           "Connect to a error database",
			dsn:            &DataSourceName{"ohyee", "123456", "127.0.0.2", 3307},
			isClose:        true,
			wantErr:        false,
			wantCloseError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConn, err := NewConnection(tt.dsn)
			if gotConn.isClose != tt.isClose || (err != nil) != tt.wantErr {
				t.Errorf("NewConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = gotConn.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect.Close() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestConnect_Query(t *testing.T) {
	if !isLocal {
		return
	}
	conn, err := NewConnection(&validDatabase)
	if err != nil {
		t.Fail()
	}
	defer conn.Close()
	user := NewUser()
	timeZone, err := time.LoadLocation("UTC")

	type args struct {
		sqlStr string
		args   []any
	}
	tests := []struct {
		name       string
		args       args
		wantResult Results
		wantErr    bool
	}{
		{
			name:       "test table",
			args:       args{"SELECT * FROM `test`.`test`;", []any{}},
			wantResult: Results{{"name": "test"}},
			wantErr:    !isLocal,
		},
		{
			name: "test common types",
			args: args{user.SQL().Query(), []any{}},
			wantResult: Results{{
				"username": "OhYee",
				"fullname": "名字",
				"password": "123123",
				"bigint":   int64(1111111111111111),
				"longtext": "qwe",
				"bit":      byte(1),
				"date":     time.Date(2019, 01, 01, 0, 0, 0, 0, timeZone),
				"datetime": time.Date(2019, 01, 01, 10, 10, 10, 0, timeZone),
				"decimal":  6666.666667,
				"double":   1.2345678910111213,
				"float":    1.23457,
			}},
			wantErr: !isLocal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := conn.Query(tt.args.sqlStr, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !goutils.Equal(gotResult, tt.wantResult) {
				t.Errorf("Connect.Query() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
