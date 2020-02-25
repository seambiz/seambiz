package stime_test

import (
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
			name: "date only",
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
