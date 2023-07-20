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
	s := &SQLStatement{}
	s.buffer.Reset()

	return s
}

// String returns a string representation
func (s *SQLStatement) String() string {
	return s.buffer.String()
}

// Query return SQL Statement as string und return the buffer to the pool.
func (s *SQLStatement) Query() string {
	return s.buffer.String()
}

// append a string to the sql statement and depending on @whitespace inserts a blank at the end
func (s *SQLStatement) append(whitespace bool, values ...interface{}) *SQLStatement {
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
			_, err := s.buffer.Write([]byte(" "))
			if err != nil {
				panic(err)
			}
		}
	}
	return s
}

// Append a string to the sql statement and a space at the end
func (s *SQLStatement) Append(values ...interface{}) *SQLStatement {
	return s.append(true, values...)
}

// AppendRaw a string to the sql statement and a space at the end
func (s *SQLStatement) AppendRaw(values ...interface{}) *SQLStatement {
	return s.append(false, values...)
}

// AppendStr a string to the sql statement and a space at the end
func (s *SQLStatement) AppendStr(strs ...string) *SQLStatement {
	for _, str := range strs {
		_, err := s.buffer.WriteString(str)
		if err != nil {
			panic(err)
		}
	}
	return s
}

func (s *SQLStatement) InInt(ints []int) *SQLStatement {
	if ints == nil {
		return s
	}

	for i, v := range ints {
		if i > 0 {
			s.AppendStr(",")
		}
		s.AppendInt(v)
	}

	return s
}

// AppendStr a string to the sql statement and a space at the end
func (s *SQLStatement) AppendStrs(prefix string, suffix string, strs ...string) *SQLStatement {
	for _, str := range strs {
		s.AppendStr(prefix, str, suffix)
	}
	return s
}

// AppendStr a string to the sql statement and a space at the end
func (s *SQLStatement) AppendBytes(whitespace bool, bs ...[]byte) *SQLStatement {
	for _, b := range bs {
		_, err := s.buffer.Write(b)
		if err != nil {
			panic(err)
		}

		if whitespace {
			_, err := s.buffer.Write([]byte(" "))
			if err != nil {
				panic(err)
			}
		}
	}
	return s
}

// appendUInt appends a string to the sql statement
func (s *SQLStatement) appendUInt(n uint) {
	s.buffer.Write(strconv.AppendInt(nil, int64(n), 10))
}

// AppendInt appends a string to the sql statement
func (s *SQLStatement) AppendInt(n int) *SQLStatement {
	s.buffer.Write(strconv.AppendInt(nil, int64(n), 10))
	return s
}

// Reset the underlying buffer.
func (s *SQLStatement) Reset() {
	s.buffer.Reset()
}

// Fields appends alle fields from a struct
func (s *SQLStatement) Fields(prepend string, prefix string, fields []string) {
	if len(fields) > 0 {
		var p string

		if prefix != "" {
			p = prefix + "."
		}

		if prepend != "" {
			s.buffer.WriteString(prepend)
		}

		for i, f := range fields {
			if i > 0 {
				s.buffer.WriteString(", ")
			}

			s.buffer.WriteString(p)
			s.buffer.WriteString(f)
		}
	}
}

// AppendFields helper for adding fields so a select statement.
func (s *SQLStatement) AppendFields(prepend string, prefix string, separator string, append string, fields []string) {
	s.buffer.WriteString(prepend)

	for i, f := range fields {
		if i > 0 {
			s.buffer.WriteString(separator)
		}

		s.buffer.WriteString(prefix)
		s.buffer.WriteString(f)
	}

	s.buffer.WriteString(append)
}

// AppendFiller helper for adding placeholder to a insert statement.
func (s *SQLStatement) AppendFiller(prepend string, separator string, append string, filler string, n int) {
	if prepend != "" {
		s.buffer.WriteString(prepend)
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			s.buffer.WriteString(separator)
		}

		s.AppendStr(filler)
	}

	s.buffer.WriteString(append)
}
