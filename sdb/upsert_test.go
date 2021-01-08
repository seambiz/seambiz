package sdb

import "testing"

type PointerData struct {
	String    string
	StringPtr *string
	Int       int
	IntPtr    *int
}

func TestUpsertStatement_Record(t *testing.T) {
	type fields struct {
		columns []string
	}
	type args struct {
		values interface{}
	}
	var i = 5
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "start simple",
			fields: fields{
				columns: []string{"String", "StringPtr", "Int", "IntPtr"},
			},
			args: args{
				values: PointerData{String: "test"},
			},
		},
		{
			name: "int ptr",
			fields: fields{
				columns: []string{"String", "StringPtr", "Int", "IntPtr"},
			},
			args: args{
				values: PointerData{String: "test", IntPtr: &i},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UpsertStatement{
				sql:     NewSQLStatement(),
				columns: tt.fields.columns,
			}
			u.Record(tt.args.values)
		})
	}
}
