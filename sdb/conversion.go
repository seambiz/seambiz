package sdb

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/rs/zerolog/log"
)

// TimeFormat Sandard MySQL datetime format
const (
	TimeFormat = time.RFC3339
)

// ToUnsafeString converts b to string without memory allocations.
//
// The returned string is valid only until b is reachable and unmodified.
func ToUnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ParseUint parses uint from buf.
func ParseUint(buf []byte) (int, error) {
	v, n, err := parseUintBuf(buf)
	if n != len(buf) {
		return -1, errUnexpectedTrailingChar
	}
	return v, err
}

var (
	errEmptyInt               = errors.New("empty integer")
	errUnexpectedFirstChar    = errors.New("unexpected first char found. Expecting 0-9")
	errUnexpectedTrailingChar = errors.New("unexpected trailing char found. Expecting 0-9")
	errTooLongInt             = errors.New("too long int")
)

func parseUintBuf(b []byte) (int, int, error) {
	n := len(b)
	if n == 0 {
		return -1, 0, errEmptyInt
	}
	v := 0
	for i := 0; i < n; i++ {
		c := b[i]
		k := c - '0'
		if k > 9 {
			if i == 0 {
				return -1, i, errUnexpectedFirstChar
			}
			return v, i, nil
		}
		// Test for overflow.
		if v*10 < v {
			return -1, i, errTooLongInt
		}
		v = 10*v + int(k)
	}
	return v, n, nil
}

// ToInt conversion from sql.RawBytes
func ToInt(b []byte) int {
	if b == nil {
		return 0
	}
	i, err := strconv.Atoi(ToUnsafeString(b))
	if err != nil {
		log.Error().Bytes("b", b).Msg("Atoi")
		return 0
	}
	return i
}

// ToBool conversion from sql.RawBytes
func ToBool(b []byte) bool {
	if b == nil {
		return false
	}
	return ToInt(b) == 1
}

// ToString conversion from sql.RawBytes
// toUnsafeString is not used, because of the limited validity of the raw bytes.
func ToString(b []byte) string {
	if b == nil {
		return ""
	}
	return string(b)
}

// ToInt64 conversion from sql.RawBytes
func ToInt64(b []byte) int64 {
	if b == nil {
		return 0
	}
	i, err := strconv.ParseInt(ToUnsafeString(b), 10, 64)
	if err != nil {
		log.Error().Bytes("b", b).Msg("ToInt64")
		return 0
	}
	return i
}

// ToUInt conversion from sql.RawBytes
func ToUInt(b []byte) uint {
	if b == nil {
		return 0
	}
	i, err := strconv.ParseUint(ToUnsafeString(b), 10, 32)
	if err != nil {
		log.Error().Bytes("b", b).Msg("ToUInt")
		return 0
	}
	return uint(i)
}

// ToUInt64 conversion from sql.RawBytes
func ToUInt64(b []byte) uint64 {
	if b == nil {
		return 0
	}
	i, err := strconv.ParseUint(ToUnsafeString(b), 10, 64)
	if err != nil {
		log.Error().Bytes("b", b).Msg("ToUInt64")
		return 0
	}
	return i
}

// ToFloat32 conversion from sql.RawBytes
func ToFloat32(b []byte) float32 {
	if b == nil {
		return 0
	}
	f, err := strconv.ParseFloat(ToUnsafeString(b), 32)
	if err != nil {
		log.Error().Bytes("b", b).Msg("ToFloat32")
		return 0
	}
	return float32(f)
}

// ToFloat64 conversion from sql.RawBytes
func ToFloat64(b []byte) float64 {
	if b == nil {
		return 0
	}
	f, err := strconv.ParseFloat(ToUnsafeString(b), 64)
	if err != nil {
		log.Error().Bytes("b", b).Msg("ToFloat64")
		return 0
	}
	return f
}

// ToTime conversion from sql.RawBytes
func ToTime(b []byte) time.Time {
	if b == nil {
		return time.Time{}
	}
	format := TimeFormat[:19]
	s := ToUnsafeString(b)
	if strings.Contains(s, "Z") {
		format = time.RFC3339
	}
	switch len(s) {
	case 8:
		if s == "00:00:00" {
			return time.Time{}
		}
		format = format[11:19]
	case 10:
		if s == "0000-00-00" {
			return time.Time{}
		}
		format = format[:10]
	case 19:
		if s == "0000-00-00 00:00:00" {
			return time.Time{}
		}
	}
	t, err := time.ParseInLocation(format, s, time.Local)
	if err != nil {
		log.Error().Bytes("b", b).Msg("ToTime")
		return time.Time{}
	}
	return t
}
