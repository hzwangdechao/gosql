package sql

// Type of database value
type Type uint8

// Database value type list
const (
	ColumnTypeFunction Type = iota
	ColumnTypeUnknown
)

var typeName = [...]string{"Function", "Unknown"}

// String name of database value type
func (t Type) String() string {
	if int(t) > len(typeName) {
		return "Unknown"
	}
	return typeName[t]
}