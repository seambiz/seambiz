package sdb

import (
	"reflect"
	"strconv"

	"unsafe"

	"time"

	"github.com/seambiz/seambiz/stime"
	"github.com/valyala/bytebufferpool"
)

// JsonBuffer type alias for shorter template functions
type JsonBuffer bytebufferpool.ByteBuffer

// +gochecknoglobals
var (
	JSONQuote             = []byte(`"`)
	strBackslashQuote     = []byte(`\"`)
	strBackslashBackslash = []byte(`\\`)
	strBackslashN         = []byte(`\n`)
	strBackslashR         = []byte(`\r`)
	strBackslashT         = []byte(`\t`)
	strBackslashF         = []byte(`\u000c`)
	strBackslashB         = []byte(`\u0008`)
	strBackslashLT        = []byte(`\u003c`)
	strBackslashQ         = []byte(`\u0027`)
	strBackslashZero      = []byte(`\u0000`)
)

// NewJsonateBuffer factory
func NewJsonateBuffer() *JsonBuffer {
	return (*JsonBuffer)(bytebufferpool.Get())
}

// NewLine write newline char to buffer
func (t *JsonBuffer) NewLine() {
	_, err := bb(t).Write([]byte{'\n'})
	if err != nil {
		panic(err)
	}
}

// SS shortcut for writing string to buffer and check error
func (t *JsonBuffer) SS(ss ...string) {
	for _, s := range ss {
		t.S(s)
	}
}

// S shortcut for writing string to buffer and check error
func (t *JsonBuffer) S(s string) {
	if s != "" {
		_, err := bb(t).WriteString(s)
		if err != nil {
			panic(err)
		}
	}
}

// JSe shortcut for writing string to buffer and check error
func (t *JsonBuffer) JSe(s string) {
	t.S(`"`)
	if s != "" {
		t.jsonString(s)
	}
	t.S(`"`)
}

// JB shortcut for writing string to JSON escaped string
func (t *JsonBuffer) JB(prepend, key string, value bool) {
	if prepend != "" {
		t.S(prepend)
	}
	bb(t).Write(JSONQuote)
	t.S(key)
	if value {
		t.S(`":true`)

	} else {
		t.S(`":false`)

	}
}

// JT shortcut for writing string to JSON escaped string
func (t *JsonBuffer) JT(prepend, key string, value time.Time) {
	if prepend != "" {
		t.S(prepend)
	}
	bb(t).Write(JSONQuote)
	t.S(key)
	t.S(`":"`)
	t.S(stime.FormatTime(&value))
	bb(t).Write(JSONQuote)
}

// JS shortcut for writing string to JSON escaped string
func (t *JsonBuffer) JS(prepend, key, value string) {
	if prepend != "" {
		t.S(prepend)
	}
	bb(t).Write(JSONQuote)
	t.S(key)
	t.S(`":"`)
	if value != "" {
		t.jsonString(value)
	}
	bb(t).Write(JSONQuote)
}

// JDu shortcut for writing int to JSON escaped string
func (t *JsonBuffer) JDu(prepend, key string, value uint) {
	if prepend != "" {
		t.S(prepend)
	}
	bb(t).Write(JSONQuote)
	t.S(key)
	t.S(`":`)
	t.Du(value)
}

// JByte shortcut for writing int to JSON escaped string
func (t *JsonBuffer) JByte(prepend, key string, value []byte) {
	if prepend != "" {
		t.S(prepend)
	}
	bb(t).Write(JSONQuote)
	t.S(key)
	t.S(`":"`)
	bb(t).Write(value)
	t.S(`"`)
}

// JF shortcut for writing int to JSON escaped string
func (t *JsonBuffer) JF(prepend, key string, value float32) {
	if prepend != "" {
		t.S(prepend)
	}
	bb(t).Write(JSONQuote)
	t.S(key)
	t.S(`":`)
	t.F(value)
}

// JF shortcut for writing int to JSON escaped string
func (t *JsonBuffer) JF64(prepend, key string, value float64) {
	if prepend != "" {
		t.S(prepend)
	}
	bb(t).Write(JSONQuote)
	t.S(key)
	t.S(`":`)
	t.F64(value)
}

// JD shortcut for writing int to JSON escaped string
func (t *JsonBuffer) JD(prepend, key string, value int) {
	if prepend != "" {
		t.S(prepend)
	}
	bb(t).Write(JSONQuote)
	t.S(key)
	t.S(`":`)
	t.D(value)
}

// JD shortcut for writing int to JSON escaped string
func (t *JsonBuffer) JD64(prepend, key string, value int64) {
	if prepend != "" {
		t.S(prepend)
	}
	bb(t).Write(JSONQuote)
	t.S(key)
	t.S(`":`)
	t.D64(value)
}

// JD shortcut for writing int to JSON escaped string
func (t *JsonBuffer) JD64u(prepend, key string, value uint64) {
	if prepend != "" {
		t.S(prepend)
	}
	bb(t).Write(JSONQuote)
	t.S(key)
	t.S(`":`)
	t.D64u(value)
}

// only for testing
func (t *JsonBuffer) Reset() {
	bb(t).Reset()
}

// space inserts space
func (t *JsonBuffer) Space() {
	t.S(" ")
}

// F append integer without allocation
func (t *JsonBuffer) F(f float32) {
	bb(t).B = strconv.AppendFloat(bb(t).B, float64(f), 'f', 2, 64)
}

// F64 append integer without allocation
func (t *JsonBuffer) F64(f float64) {
	bb(t).B = strconv.AppendFloat(bb(t).B, f, 'f', 2, 64)
}

// D append integer without allocation
func (t *JsonBuffer) D(n int) {
	bb(t).B = strconv.AppendInt(bb(t).B, int64(n), 10)
}

// Du append integer without allocation
func (t *JsonBuffer) Du(n uint) {
	bb(t).B = strconv.AppendUint(bb(t).B, uint64(n), 10)
}

// D64u append integer without allocation
func (t *JsonBuffer) D64u(n uint64) {
	bb(t).B = strconv.AppendUint(bb(t).B, n, 10)
}

// D64 append integer without allocation
func (t *JsonBuffer) D64(n int64) {
	bb(t).B = strconv.AppendInt(bb(t).B, n, 10)
}

// Bytes returns buffer contents
// Attention: returns the buffer to the pool and sets pointer to nil
func (t *JsonBuffer) Bytes() []byte {
	b := bb(t).Bytes()
	bytebufferpool.Put(bb(t))
	return b
}

func bb(t *JsonBuffer) *bytebufferpool.ByteBuffer {
	return (*bytebufferpool.ByteBuffer)(t)
}

func (t *JsonBuffer) jsonString(s string) {
	write := bb(t).Write
	b := s2b(s)
	j := 0
	n := len(b)
	if n > 0 {
		// Hint the compiler to remove bounds checks in the loop below.
		_ = b[n-1]
	}
	for i := 0; i < n; i++ {
		switch b[i] {
		case '"':
			write(b[j:i])
			write(strBackslashQuote)
			j = i + 1
		case '\\':
			write(b[j:i])
			write(strBackslashBackslash)
			j = i + 1
		case '\n':
			write(b[j:i])
			write(strBackslashN)
			j = i + 1
		case '\r':
			write(b[j:i])
			write(strBackslashR)
			j = i + 1
		case '\t':
			write(b[j:i])
			write(strBackslashT)
			j = i + 1
		case '\f':
			write(b[j:i])
			write(strBackslashF)
			j = i + 1
		case '\b':
			write(b[j:i])
			write(strBackslashB)
			j = i + 1
		case '<':
			write(b[j:i])
			write(strBackslashLT)
			j = i + 1
		case '\'':
			write(b[j:i])
			write(strBackslashQ)
			j = i + 1
		case 0:
			write(b[j:i])
			write(strBackslashZero)
			j = i + 1
		}
	}
	write(b[j:])
}

// b2s converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
