package stime

import "time"

// Now return unix timestamp as uint32
func Now() uint {
	t := time.Now().Unix()
	return uint(t)
}

// endOfDay helper copied from jinzhu/now
func endOfDay(t *time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// ParseDate parses german date format and returns unix timestamp
func ParseDate(datum string, endofday bool) (uint, error) {
	t, err := time.Parse("02.01.2006", datum)
	if err != nil {
		return 0, err
	}

	if endofday {
		return uint(endOfDay(&t).Unix()), nil
	}

	return uint(t.Unix()), nil
}

// ParseDateOnly returns date part as integer
func ParseDateOnly(datum string) (int, error) {
	t, err := time.Parse("02.01.2006", datum)
	if err != nil {
		return 0, err
	}

	return t.Year()*10000 + int(t.Month()*100) + t.Day(), nil
}

// Format unix timestamp in german date notation
func Format(t uint) string {
	date := time.Unix(int64(t), 0)
	return date.Format("02.01.2006")
}

// FormatFull unix timestamp into german date and time notation
func FormatFull(t uint) string {
	date := time.Unix(int64(t), 0)
	return date.Format("02.01.2006 15:04:05")
}

// FormatTime time.Time to german date notation
func FormatTime(t *time.Time) string {
	return t.Format("02.01.2006")
}
