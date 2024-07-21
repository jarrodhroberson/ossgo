package timestamp

import (
	"reflect"
	"testing"
	"time"
)

func TestToday(t *testing.T) {
	now := Now()
	tests := []struct {
		name string
		want Period
	}{
		{
			name: "test today",
			want: Period{
				Start: now.ZeroTime(),
				End:   now.Add(time.Hour * 24).ZeroTime(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Today(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Today() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addMonth(t *testing.T) {
	type args struct {
		t time.Time
		m int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "test add 1 month from October 31",
			args: args{
				t: time.Date(2024, time.October, 31, 1, 2, 3, 0, time.UTC),
				m: 1,
			},
			want: time.Date(2024, time.November, 30, 1, 2, 3, 0, time.UTC),
		},
		{
			name: "test add 1 month from December 31, 2024",
			args: args{
				t: time.Date(2024, time.December, 31, 1, 2, 3, 0, time.UTC),
				m: 1,
			},
			want: time.Date(2025, time.January, 31, 1, 2, 3, 0, time.UTC),
		},
		{
			name: "test add 1 month from October 31",
			args: args{
				t: time.Date(2024, time.October, 31, 1, 2, 3, 0, time.UTC),
				m: 1,
			},
			want: time.Date(2024, time.November, 30, 1, 2, 3, 0, time.UTC),
		},
		{
			name: "test add 1 month from February 29",
			args: args{
				t: time.Date(2024, time.February, 29, 1, 2, 3, 0, time.UTC),
				m: 1,
			},
			want: time.Date(2024, time.March, 29, 1, 2, 3, 0, time.UTC),
		},
		{
			name: "test add 1 month from January 29 on a leap year",
			args: args{
				t: time.Date(2024, time.January, 29, 1, 2, 3, 0, time.UTC),
				m: 1,
			},
			want: time.Date(2024, time.February, 29, 1, 2, 3, 0, time.UTC),
		},
		{
			name: "test add 1 month from January 30 on a leap year",
			args: args{
				t: time.Date(2024, time.January, 30, 1, 2, 3, 0, time.UTC),
				m: 1,
			},
			want: time.Date(2024, time.February, 29, 1, 2, 3, 0, time.UTC),
		},
		{
			name: "test add 1 month from January 31 on a leap year",
			args: args{
				t: time.Date(2024, time.January, 31, 1, 2, 3, 0, time.UTC),
				m: 1,
			},
			want: time.Date(2024, time.February, 29, 1, 2, 3, 0, time.UTC),
		},
		{
			name: "test add 1 month from January 31 on a non-leap year",
			args: args{
				t: time.Date(2024, time.January, 31, 1, 2, 3, 0, time.UTC),
				m: 1,
			},
			want: time.Date(2024, time.February, 29, 1, 2, 3, 0, time.UTC),
		},
		{
			name: "test add 1 month from January 30 on a non-leap year",
			args: args{
				t: time.Date(2023, time.January, 30, 1, 2, 3, 0, time.UTC),
				m: 1,
			},
			want: time.Date(2023, time.February, 28, 1, 2, 3, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addMonth(tt.args.t, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToRange(t *testing.T) {
	type args struct {
		from     Timestamp
		to       Timestamp
		interval int64
		d        time.Duration
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Hours of Day",
			args: args{
				from: Timestamp{
					t: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				to: Timestamp{
					t: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
				},
				interval: 24,
				d:        time.Hour * 24,
			},
			want: 24,
		},
		{
			name: "Days of Month",
			args: args{
				from: Timestamp{
					t: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				to: Timestamp{
					t: time.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC),
				},
				interval: 31,
				d:        (time.Hour * 24) * 7,
			},
			want: 31,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToRange(tt.args.from, tt.args.to, tt.args.interval, tt.args.d); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("ToRange() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
