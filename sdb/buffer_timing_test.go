package sdb

import (
	"testing"
)

func BenchmarkAppend(b *testing.B) {
	s := "foobarbaz"
	b.RunParallel(func(pb *testing.PB) {
		sql := NewSQLStatement()
		for pb.Next() {
			for i := 0; i < 100; i++ {
				sql.Append(s, s)
			}
			sql.Reset()
		}
	})
}

func BenchmarkAppendInt(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		sql := NewSQLStatement()
		for pb.Next() {
			sql.AppendRaw(1521)

			sql.Reset()
		}
	})
}

func BenchmarkAppendStr(b *testing.B) {
	s := "foobarbaz"
	b.RunParallel(func(pb *testing.PB) {
		sql := NewSQLStatement()
		for pb.Next() {
			for i := 0; i < 100; i++ {
				sql.AppendStr(s, s)
			}
			sql.Reset()
		}
	})
}

func BenchmarkAppendBytes(b *testing.B) {
	s := []byte("foobarbaz")
	b.RunParallel(func(pb *testing.PB) {
		sql := NewSQLStatement()
		for pb.Next() {
			for i := 0; i < 100; i++ {
				sql.AppendBytes(false, s, s)
			}
			sql.Reset()
		}
	})
}

func BenchmarkAppendFields(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		sql := NewSQLStatement()
		for pb.Next() {
			for i := 0; i < 100; i++ {
				sql.AppendFields("SELECT ", "A.", ",", " FROM table A", []string{"f1", "f2", "f3"})
			}
			sql.Reset()
		}
	})
}

func BenchmarkAppendFiller(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		sql := NewSQLStatement()
		for pb.Next() {
			for i := 0; i < 100; i++ {
				sql.AppendFiller("(", ",", ")", "?", 5)
			}
			sql.Reset()
		}
	})
}
