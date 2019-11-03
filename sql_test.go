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

	tb := NewTable("test", "user", as[0])
	return &User{
		// SQL:      NewSQL().From(tb),
		Table:    tb,
		Username: NewColumn(tb, "username", as[1], ColumnTypeUnknown),
		Realname: NewColumn(tb, "realname", as[2], ColumnTypeUnknown),
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
		{
			name:  "select max",
			query: user.SQL().Select(user.Username, NewFunctionColumn("MAX(%s)", "max", user.Realname.String())),
			sql:   "SELECT `test`.`user`.`username`, MAX(`test`.`user`.`realname`) AS `max` FROM `test`.`user`;",
		},
		{
			name:  "select max without as",
			query: user.SQL().Select(user.Username, NewFunctionColumn("MAX(%s)", "", user.Realname.String())),
			sql:   "SELECT `test`.`user`.`username`, MAX(`test`.`user`.`realname`) FROM `test`.`user`;",
		},
		{
			name:  "from table alias",
			query: user.As("usr").SQL().Select(NewUser("usr", "name").Username),
			sql:   "SELECT `usr`.`username` AS `name` FROM `test`.`user` AS `usr`;",
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
