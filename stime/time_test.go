package stime_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/seambiz/seambiz/stime"
)

func TestNow(t *testing.T) {
	tests := []struct {
		name string
		want uint
	}{
		{
			name: "now",
			want: uint(time.Now().Unix()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stime.Now(); got != tt.want {
				t.Errorf("Now() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	type args struct {
		datum    string
		endofday bool
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{
			name: "date only",
			args: args{
				datum:    "10.12.2016",
				endofday: false,
			},
			want:    1481328000,
			wantErr: false,
		},
		{
			name: "date endofday",
			args: args{
				datum:    "10.12.2016",
				endofday: true,
			},
			want:    1481414399,
			wantErr: false,
		},
		{
			name: "format error",
			args: args{
				datum:    "99.99.2016",
				endofday: true,
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stime.ParseDate(tt.args.datum, tt.args.endofday)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDateOnly(t *testing.T) {
	type args struct {
		datum string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "date only",
			args: args{
				datum: "10.12.2016",
			},
			want:    20161210,
			wantErr: false,
		},
		{
			name: "format error",
			args: args{
				datum: "99.99.2016",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stime.ParseDateOnly(tt.args.datum)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateOnly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDateOnly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	type args struct {
		format string
		t      uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "date only",
			args: args{
				format: "02.01.2006",
				t:      1481328000,
			},
			want: "10.12.2016",
		},
		{
			name: "rfc3339",
			args: args{
				format: time.RFC3339,
				t:      1481414399,
			},
			want: "2016-12-11T00:59:59+01:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stime.Format(tt.args.format, tt.args.t); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	type args struct {
		t uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "date only",
			args: args{
				t: 1481328000,
			},
			want: "10.12.2016",
		},
		{
			name: "date only",
			args: args{
				t: 1481414399,
			},
			want: "11.12.2016",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stime.FormatDate(tt.args.t); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatFull(t *testing.T) {
	type args struct {
		t uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "date only",
			args: args{
				t: 1481328000,
			},
			want: "10.12.2016 01:00:00",
		},
		{
			name: "date only",
			args: args{
				t: 1481414399,
			},
			want: "11.12.2016 00:59:59",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stime.FormatFull(tt.args.t); got != tt.want {
				t.Errorf("FormatFull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOtherTimezone(t *testing.T) {
	type args struct {
		t      uint
		tz     string
		format string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "los angeles",
			args: args{
				t:  1481328000,
				tz: "America/Los_Angeles",
			},
			want: "2016-12-09T16:00:00-08:00",
		},
		{
			name: "hongkong",
			args: args{
				t:  1481328000,
				tz: "Asia/Hong_Kong",
			},
			want: "2016-12-10T08:00:00+08:00",
		},
		{
			name: "date only",
			args: args{
				t:      1481328000,
				tz:     "Asia/Hong_Kong",
				format: "02.01.2006",
			},
			want: "10.12.2016",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stime.FormatIn(tt.args.t, tt.args.tz, tt.args.format); got != tt.want {
				t.Errorf("FormatIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "date only",
			args: args{
				t: time.Date(2013, 11, 18, 17, 51, 49, 123456789, time.UTC),
			},
			want: "18.11.2013",
		},
		{
			name: "date only",
			args: args{
				t: time.Date(2016, 02, 11, 0, 0, 0, 123456789, time.UTC),
			},
			want: "11.02.2016",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stime.FormatTime(&tt.args.t); got != tt.want {
				t.Errorf("FormatTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIn(t *testing.T) {
	europeBerlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Error("timezone could not be loaded")
	}
	americaLosAngeles, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Error("timezone could not be loaded")
	}

	type args struct {
		unix    uint
		locIANA string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "empty loc",
			args: args{
				unix:    uint(time.Date(2013, 11, 18, 17, 51, 49, 0, time.UTC).Unix()),
				locIANA: "",
			},
			want: time.Date(2013, 11, 18, 17, 51, 49, 0, time.UTC),
		},
		{
			name: "europe",
			args: args{
				unix:    uint(time.Date(2013, 11, 18, 17, 51, 49, 0, time.UTC).Unix()),
				locIANA: "Europe/Berlin",
			},
			want: time.Date(2013, 11, 18, 18, 51, 49, 0, europeBerlin),
		},
		{
			name: "europe",
			args: args{
				unix:    uint(time.Date(2013, 11, 18, 17, 51, 49, 0, time.UTC).Unix()),
				locIANA: "America/Los_Angeles",
			},
			want: time.Date(2013, 11, 18, 9, 51, 49, 0, americaLosAngeles),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stime.In(tt.args.unix, tt.args.locIANA); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
		})
	}
}
