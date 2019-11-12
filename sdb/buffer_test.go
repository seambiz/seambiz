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
	want := "id, test, third "
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestFieldsFull(t *testing.T) {
	sql := sdb.NewSQLStatement()
	fields := []string{"id", "test", "third"}

	sql.Fields("'prepend', ", "abc", fields)
	got := sql.Query()
	want := "'prepend', abc.id, abc.test, abc.third "
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
	want := "a.id, a.test, a.third ,b.id, b.test, b.third "
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
