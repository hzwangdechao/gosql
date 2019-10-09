package sql

import (
	"testing"
)

type User struct {
	// *SQL
	*Table
	Username *Column
	Realname *Column
}

func NewUser() *User {
	tb := NewTable("test", "user", "")
	return &User{
		// SQL:      NewSQL().From(tb),
		Table:    tb,
		Username: NewColumn(tb, "username", "", ColumnTypeUnknow),
		Realname: NewColumn(tb, "realname", "", ColumnTypeUnknow),
	}
}

func Test_main(t *testing.T) {
	user := NewUser()
	testcases := []struct {
		name  string
		query *SQL
		sql   string
	}{
		{
			name:  "select from",
			query: user.SQL().Select(user.Username, user.Realname),
			sql:   "SELECT `test`.`user`.`username`, `test`.`user`.`realname` FROM `test`.`user`;",
		},
		{
			name:  "select * from",
			query: user.SQL(),
			sql:   "SELECT * FROM `test`.`user`;",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			if q := tt.query.Query(); q != tt.sql {
				t.Errorf("want %s, got %s\n", tt.sql, q)
			}
		})
	}
	//
	// fmt.Println(user.Select(user.Username, user.Realname).Query())
}
