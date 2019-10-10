package op

import (
	"github.com/OhYee/gosql"
	"testing"
)

func Test_expression(t *testing.T) {
	testcases := []struct {
		name string
		op   *Operator
		want string
	}{
		{
			name: "Eq",
			op:   Eq(1, 2),
			want: "(1 = 2)",
		},
		{
			name: "Ne",
			op:   Ne(sql.NewColumn(sql.NewTable("test", "user", ""), "username", "", sql.ColumnTypeUnknow), "OhYee"),
			want: "(`test`.`user`.`username` != 'OhYee')",
		},
		{
			name: "Lt",
			op:   Lt("`test`.`user`.`username`", "OhYee"),
			want: "(`test`.`user`.`username` < 'OhYee')",
		},
		{
			name: "Gt",
			op:   Gt("`test`.`user`.`username`", "'OhYee'"),
			want: "(`test`.`user`.`username` > 'OhYee')",
		},
		{
			name: "Ge",
			op:   Ge("苟利国家生死以", "'岂因祸福避趋之'"),
			want: "('苟利国家生死以' >= '岂因祸福避趋之')",
		},
		{
			name: "Le",
			op:   Le("`age`", 15),
			want: "(`age` <= 15)",
		},
		{
			name: "Like",
			op:   Like("`words`", "%测试%"),
			want: "(`words` LIKE '%测试%')",
		},
		{
			name: "In",
			op:   In("`username`", []string{"OhYee", "Abcdefg", "oyohyee@oyohyee.com"}),
			want: "(`username` IN ('OhYee', 'Abcdefg', 'oyohyee@oyohyee.com'))",
		},
		{
			name: "Not in",
			op:   Ni("`age`", []int{15, 16, 17}),
			want: "(`age` NOT IN (15, 16, 17))",
		},
		{
			name: "Between",
			op:   Between("`age`", 15, 17),
			want: "(`age` BETWEEN 15 AND 17)",
		},
		{
			name: "Not",
			op:   Not(In("`age`", []int{15, 16, 17})),
			want: "(NOT (`age` IN (15, 16, 17)))",
		},
		{
			name: "Other string",
			op:   Str("a > 1"),
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
