package sdb

import (
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
var sqlBuffer = &sync.Pool{
	New: func() interface{} {
		return &SQLStatement{
			buffer: make([]byte, 0, 1024),
		}
	},
}

// SQLStatement is a wrapper around bytepufferpool for nicer usage
type SQLStatement struct {
	buffer       []byte
	fieldsCalled bool
}

// NewSQLStatement return bytebuffer for a statement
func NewSQLStatement() *SQLStatement {
	return sqlBuffer.Get().(*SQLStatement)
}

// String returns a string representation
func (s *SQLStatement) String() string {
	return string(s.buffer)
}

// Query return SQL Statement as string und return the buffer to the pool.
func (s *SQLStatement) Query() string {
	defer sqlBuffer.Put(s)
	defer s.Reset()

	return s.String()
}

func (s *SQLStatement) Bytes() []byte {
	defer sqlBuffer.Put(s)
	defer s.Reset()

	return s.buffer
}

// append a string to the sql statement and depending on @whitespace inserts a blank at the end
func (s *SQLStatement) append(whitespace bool, values ...interface{}) *SQLStatement {
	for _, v := range values {
		switch v := v.(type) {
		case string:
			_, err := s.WriteString(v)
			if err != nil {
				panic(err)
			}

		case int:
			s.AppendInt(v)
		case uint:
			s.appendUInt(v)
		}

		if whitespace {
			_, err := s.Write([]byte(" "))
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
		_, err := s.WriteString(str)
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
		_, err := s.Write(b)
		if err != nil {
			panic(err)
		}

		if whitespace {
			_, err := s.Write([]byte(" "))
			if err != nil {
				panic(err)
			}
		}
	}
	return s
}

// appendUInt appends a string to the sql statement
func (s *SQLStatement) appendUInt(n uint) {
	s.Write(strconv.AppendInt(nil, int64(n), 10))
}

// AppendInt appends a string to the sql statement
func (s *SQLStatement) AppendInt(n int) *SQLStatement {
	s.buffer = strconv.AppendInt(s.buffer, int64(n), 10)

	return s
}

// Reset the underlying buffer.
func (s *SQLStatement) Reset() {
	s.buffer = s.buffer[:0]
	s.fieldsCalled = false
}

// Write implements io.Writer - it appends p to ByteBuffer.B
func (s *SQLStatement) Write(p []byte) (int, error) {
	s.buffer = append(s.buffer, p...)
	return len(p), nil
}

// WriteString appends s to ByteBuffer.B.
func (s *SQLStatement) WriteString(str string) (int, error) {
	s.buffer = append(s.buffer, str...)
	return len(str), nil
}

// Fields appends alle fields from a struct
func (s *SQLStatement) Fields(prefix string, fields []string) {
	if len(fields) > 0 {
		if s.fieldsCalled {
			s.WriteString(",")
		}
		s.fieldsCalled = true

		for i, f := range fields {
			if i > 0 {
				s.WriteString(",")
			}

			if prefix != "" {
				s.AppendStr(prefix, ".", f)
			} else {
				s.AppendStr(f)
			}
		}

		s.WriteString(" ")
	}
}

// AppendFields helper for adding fields so a select statement.
func (s *SQLStatement) AppendFields(prepend string, prefix string, separator string, append string, fields []string) {
	s.WriteString(prepend)

	for i, f := range fields {
		if i > 0 {
			s.WriteString(separator)
		}

		s.WriteString(prefix)
		s.WriteString(f)
	}

	s.WriteString(append)
}

// AppendFiller helper for adding placeholder to a insert statement.
func (s *SQLStatement) AppendFiller(prepend string, separator string, append string, filler string, n int) {
	if prepend != "" {
		s.WriteString(prepend)
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			s.WriteString(separator)
		}

		s.AppendStr(filler)
	}

	s.WriteString(append)
}
