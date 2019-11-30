package connect

import (
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var validDatabase = DataSourceName{"ohyee", "123456", "127.0.0.1", 3306}

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
			isClose:        false,
			wantErr:        false,
			wantCloseError: false,
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
	conn, err := NewConnection(&validDatabase)
	if err != nil {
		t.Fail()
	}
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
			wantResult: Results{{"name": []byte("test")}},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := conn.Query(tt.args.sqlStr, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Connect.Query() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
