package sdb

import (
	"bytes"
	"strconv"
	"sync"
)

const (
	// INNER join word
	INNER = "INNER"
	// LEFT join word
	LEFT = "LEFT"
	// RIGHT join word
	RIGHT = "RIGHT"
	// OUTER join word
	OUTER = "OUTER"
)

// nolint[gochecknoblobals]
var sqlBuffer = sync.Pool{
	New: func() interface{} {
		return &SQLStatement{}
	},
}

// SQLStatement is a wrapper around bytepufferpool for nicer usage
type SQLStatement struct {
	buffer bytes.Buffer
}

// NewSQLStatement return bytebuffer for a statement
func NewSQLStatement() *SQLStatement {
	return sqlBuffer.Get().(*SQLStatement)
	// return &SQLStatement{}
}

// String returns a string representation
func (s *SQLStatement) String() string {
	return s.buffer.String()
}

// Query return SQL Statement as string und return the buffer to the pool.
func (s *SQLStatement) Query() string {
	defer s.buffer.Reset()
	defer sqlBuffer.Put(s)
	return s.buffer.String()
}

// append a string to the sql statement and depending on @whitespace inserts a blank at the end
func (s *SQLStatement) append(whitespace bool, values ...interface{}) {
	for _, v := range values {
		switch v := v.(type) {
		case string:
			_, err := s.buffer.WriteString(v)
			if err != nil {
				panic(err)
			}

		case int:
			s.AppendInt(v)
		case uint:
			s.appendUInt(v)
		}

		if whitespace {
			_, err := s.buffer.WriteString(" ")
			if err != nil {
				panic(err)
			}
		}
	}
}

// Append a string to the sql statement and a space at the end
func (s *SQLStatement) Append(values ...interface{}) {
	s.append(true, values...)
}

// AppendRaw a string to the sql statement and a space at the end
func (s *SQLStatement) AppendRaw(values ...interface{}) {
	s.append(false, values...)
}

// appendUInt appends a string to the sql statement
func (s *SQLStatement) appendUInt(n uint) {
	s.buffer.Write(strconv.AppendInt(nil, int64(n), 10))
}

// AppendInt appends a string to the sql statement
func (s *SQLStatement) AppendInt(n int) {
	s.buffer.Write(strconv.AppendInt(nil, int64(n), 10))
}

// Fields appends alle fields from a struct
func (s *SQLStatement) Fields(prepend string, prefix string, fields []string) {
	if len(fields) > 0 {
		var p string
		if prefix != "" {
			p = prefix + "."
		}
		if prepend != "" {
			s.AppendRaw(prepend)
		}
		for i, f := range fields {
			if i > 0 {
				s.AppendRaw(`, `)
			}
			s.AppendRaw(p, f)
		}
		s.AppendRaw(" ")
	}
}
