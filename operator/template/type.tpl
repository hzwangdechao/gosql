// OperatorType expression operator type
type OperatorType uint8

// operator constant
const (
	_ OperatorType = iota
    {{range $operator := .}}Operator{{$operator.name}}
    {{end}}
)

var operatorName = [...]string{
    "Unknown", 
    {{range $operator := .}}"{{$operator.name}}",
    {{end}}
}

func (co OperatorType) String() string {
	if int(co) > len(operatorName) {
		return "Unknown"
	}
	return operatorName[co]
}
