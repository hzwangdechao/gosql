package sql

import "testing"

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want string
	}{
		{name: "TypeFunction", t: ColumnTypeFunction, want: "Function"},
		{name: "TypeUnknown", t: ColumnTypeUnknown, want: "Unknown"},
		{name: "Unknown", t: 0xff, want: "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("Type.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
