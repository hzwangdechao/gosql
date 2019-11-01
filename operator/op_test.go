package op_test

import (
	"testing"

	"github.com/OhYee/gosql"
	"github.com/OhYee/gosql/operator"
)

type User struct {
	// *SQL
	*sql.Table
	Username *sql.Column
	Realname *sql.Column
}

func NewUser() *User {
	tb := sql.NewTable("test", "user", "")
	return &User{
		// SQL:      NewSQL().From(tb),
		Table:    tb,
		Username: sql.NewColumn(tb, "username", "", sql.ColumnTypeUnknow),
		Realname: sql.NewColumn(tb, "realname", "", sql.ColumnTypeUnknow),
	}
}

func Test_expression(t *testing.T) {
	user := NewUser()
	testcases := []struct {
		name string
		op   *op.Operator
		want string
	}{
		{
			name: "Eq",
			op:   op.Eq(1, 2),
			want: "(1 = 2)",
		},
		{
			name: "Ne",
			op:   op.Ne("`username`", ""),
			want: "(`username` != '')",
		},
		{
			name: "Lt",
			op:   op.Lt("`test`.`user`.`username`", "OhYee"),
			want: "(`test`.`user`.`username` < 'OhYee')",
		},
		{
			name: "Gt",
			op:   op.Gt("`test`.`user`.`username`", "'OhYee'"),
			want: "(`test`.`user`.`username` > 'OhYee')",
		},
		{
			name: "Ge",
			op:   op.Ge("苟利国家生死以", "'岂因祸福避趋之'"),
			want: "('苟利国家生死以' >= '岂因祸福避趋之')",
		},
		{
			name: "Le",
			op:   op.Le("`age`", 15),
			want: "(`age` <= 15)",
		},
		{
			name: "Like",
			op:   op.Like("`words`", "%测试%"),
			want: "(`words` LIKE '%测试%')",
		},
		{
			name: "In",
			op:   op.In("`username`", user.SQL().Select(user.Username)),
			want: "(`username` IN (SELECT `test`.`user`.`username` FROM `test`.`user`))",
		},
		{
			name: "Not in",
			op:   op.Ni("`age`", []int{15, 16, 17}),
			want: "(`age` NOT IN (15, 16, 17))",
		},
		{
			name: "Between",
			op:   op.Between("`age`", 15, 17),
			want: "(`age` BETWEEN 15 AND 17)",
		},
		{
			name: "Not",
			op:   op.Not(op.In("`age`", []int{15, 16, 17})),
			want: "(NOT (`age` IN (15, 16, 17)))",
		},
		{
			name: "Other string",
			op:   op.Str("a > 1"),
			want: "(a > 1)",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.String(); got != tt.want {
				t.Errorf("want %s, got %s\n", tt.want, got)
			}
		})
	}
}
