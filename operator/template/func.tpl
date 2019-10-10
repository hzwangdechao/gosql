{{range $operator := .}}
{{template "block" $operator}}
{{end}}

{{define "block"}}
{{if ne .alias ""}}
// {{.alias}} for {{.name}}
var {{.alias}} = {{.name}}
{{end}}
{{if eq .num 1.0}}
// {{.name}} {{printf .format "a"}}
func {{.name}}(a any) *Operator {
	return &Operator{
		args:     []any{a},
		operator: Operator{{.name}},
        format:   "{{.format}}",
	}
}
{{end}}
{{if eq .num 2.0}}
// {{.name}} {{printf .format "a" "b"}}
func {{.name}}(a any, b any) *Operator {
	return &Operator{
		args:     []any{a, b},
		operator: Operator{{.name}},
        format:   "{{.format}}",
	}
}
{{end}}
{{if eq .num 3.0}}
// {{.name}} {{printf .format "a" "b" "c"}}
func {{.name}}(a any, b any, c any) *Operator {
	return &Operator{
		args:     []any{a, b, c},
		operator: Operator{{.name}},
        format:   "{{.format}}",
	}
}
{{end}}
{{end}}