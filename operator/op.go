package op

import (
	"fmt"
	"reflect"
	"strings"
)

//go:generate gcg ./data.json

// Operator is commonly used SQL operators
type Operator struct {
	args     []any
	operator OperatorType
	format   string
}
type any = interface{}

// String return the string that can be identify by database
func (op *Operator) String() string {
	args := make([]any, len(op.args))
	if op.operator == OperatorStr {
		args[0] = op.args[0]
	} else {
		for i, arg := range op.args {
			args[i] = toSQLString(arg)
		}
	}
	return fmt.Sprintf(op.format, args...)
}

func toSQLString(v any) (res string) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.String:
		s := v.(string)
		sl := len(s)
		if (s[0] == '\'' && s[sl-1] == '\'') ||
			(s[0] == '"' && s[sl-1] == '"') ||
			(s[0] == '`' && s[sl-1] == '`') {
			res = fmt.Sprintf("%s", v)
		} else {
			res = fmt.Sprintf("'%s'", v)
		}
	case reflect.Slice:
		s := reflect.ValueOf(v)
		ss := make([]string, s.Len())

		for i := 0; i < s.Len(); i++ {
			ss[i] = toSQLString(s.Index(i).Interface())
		}
		res = fmt.Sprintf("(%s)", strings.Join(ss, ", "))
	default:
		res = fmt.Sprintf("%v", v)
	}
	return
}
