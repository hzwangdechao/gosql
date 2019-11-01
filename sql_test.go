package sql

import (
	"github.com/OhYee/gosql/operator"
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
		{
			name:  "select from where",
			query: user.SQL().Select(user.Username, user.Realname).Where(op.Eq(user.Username, "OhYee")),
			sql:   "SELECT `test`.`user`.`username`, `test`.`user`.`realname` FROM `test`.`user` WHERE (`test`.`user`.`username` = 'OhYee');",
		},
		{
			name:  "column as",
			query: user.SQL().Select(user.Username, user.Realname.As("truename")).Where(op.Eq(user.Username, "OhYee")),
			sql:   "SELECT `test`.`user`.`username`, `test`.`user`.`realname` AS `truename` FROM `test`.`user` WHERE (`test`.`user`.`username` = 'OhYee');",
		},
		{
			name:  "select in",
			query: user.SQL().Select(user.Username, user.Realname).Where(op.In(user.Username, user.SQL().Select(user.Username).Where(op.Ne(user.Realname, "")))),
			sql:   "SELECT `test`.`user`.`username`, `test`.`user`.`realname` FROM `test`.`user` WHERE (`test`.`user`.`username` IN (SELECT `test`.`user`.`username` FROM `test`.`user` WHERE (`test`.`user`.`realname` != '')));",
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
