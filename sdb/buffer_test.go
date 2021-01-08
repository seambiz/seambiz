package sdb_test

import (
	"testing"

	"github.com/seambiz/seambiz/sdb"
)

func TestAppend(t *testing.T) {
	sql := sdb.NewSQLStatement()
	sql.Append("SELECT")

	got := sql.String()
	if got != "SELECT " {
		t.Errorf("got '%s', want '%s'", got, "SELECT ")
	}

	got = sql.Query()
	if got != "SELECT " {
		t.Errorf("got '%s', want '%s'", got, "SELECT ")
	}
}

func TestAppendRaw(t *testing.T) {
	sql := sdb.NewSQLStatement()
	sql.AppendRaw("SELECT")

	got := sql.String()
	if got != "SELECT" {
		t.Errorf("got '%s', want '%s'", got, "SELECT")
	}

	got = sql.Query()
	if got != "SELECT" {
		t.Errorf("got '%s', want '%s'", got, "SELECT")
	}
}

func TestAppendInt(t *testing.T) {
	sql := sdb.NewSQLStatement()
	sql.AppendRaw(1521)

	got := sql.Query()
	if got != "1521" {
		t.Errorf("got '%s', want '%s'", got, "1521")
	}
}

func TestFieldsSimple(t *testing.T) {
	sql := sdb.NewSQLStatement()
	fields := []string{"id", "test", "third"}

	sql.Fields("", "", fields)
	got := sql.Query()
	want := "id, test, third"
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestFieldsFull(t *testing.T) {
	sql := sdb.NewSQLStatement()
	fields := []string{"id", "test", "third"}

	sql.Fields("'prepend', ", "abc", fields)
	got := sql.Query()
	want := "'prepend', abc.id, abc.test, abc.third"
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestFieldsCodegen(t *testing.T) {
	sql := sdb.NewSQLStatement()
	fields := []string{"id", "test", "third"}
	fields2 := []string{"id", "test", "third"}

	sql.Fields("", "a", fields)
	sql.Fields(",", "b", fields2)
	got := sql.Query()
	want := "a.id, a.test, a.third,b.id, b.test, b.third"
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestSQLStatement_AppendFields(t *testing.T) {
	type args struct {
		prepend   string
		prefix    string
		separator string
		append    string
		fields    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empy",
			args: args{},
			want: "",
		},
		{
			name: "simple field list",
			args: args{
				separator: ",",
				fields:    []string{"test", "yes"},
			},
			want: "test,yes",
		},
		{
			name: "insert use case",
			args: args{
				prepend:   "(",
				prefix:    "",
				append:    ")",
				separator: ",",
				fields:    []string{"field1", "field2"},
			},
			want: "(field1,field2)",
		},
		{
			name: "select use case",
			args: args{
				prepend:   "SELECT ",
				prefix:    "A.",
				append:    " FROM table A",
				separator: ",",
				fields:    []string{"field1", "field2"},
			},
			want: "SELECT A.field1,A.field2 FROM table A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := sdb.NewSQLStatement()
			sql.AppendFields(tt.args.prepend, tt.args.prefix, tt.args.separator, tt.args.append, tt.args.fields)
			got := sql.Query()
			if got != tt.want {
				t.Errorf("got '%s', want '%s'", got, tt.want)
			}
		})
	}
}

func TestSQLStatement_AppendFiller(t *testing.T) {
	type args struct {
		prepend   string
		separator string
		append    string
		filler    string
		n         int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{},
			want: "",
		},
		{
			name: "insert use case",
			args: args{
				prepend:   "(",
				append:    ")",
				separator: ",",
				filler:    "?",
				n:         3,
			},
			want: "(?,?,?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := sdb.NewSQLStatement()
			sql.AppendFiller(tt.args.prepend, tt.args.separator, tt.args.append, tt.args.filler, tt.args.n)
			got := sql.Query()
			if got != tt.want {
				t.Errorf("got '%s', want '%s'", got, tt.want)
			}
		})
	}
}

func TestSQLStatement_InInt(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want string
	}{
		{
			name: "empty",
			args: nil,
			want: "",
		},
		{
			name: "single",
			args: []int{99},
			want: "?",
		},
		{
			name: "three",
			args: []int{1, 2, 3},
			want: "?,?,?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := sdb.NewSQLStatement()
			sql.InInt(tt.args)
			got := sql.Query()
			if got != tt.want {
				t.Errorf("got '%s', want '%s'", got, tt.want)
			}
		})
	}
}
