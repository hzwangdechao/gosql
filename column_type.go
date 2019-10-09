package sql

// Type of database value
type Type uint8

// Database value type list
const (
	ColumnTypeFunction Type = iota
	ColumnTypeUnknow
)

var typeName = [...]string{"Function", "Unknow"}

// String name of database value type
func (t Type) String() string {
	if int(t) > len(typeName) {
		return "Unknow"
	}
	return typeName[t]
}