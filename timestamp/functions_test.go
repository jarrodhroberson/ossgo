package timestamp

import (
	"fmt"
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
		from *Timestamp
		to   *Timestamp
		d    time.Duration
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Hours of Day",
			args: args{
				from: &Timestamp{
					t: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				to: &Timestamp{
					t: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
				},
				d: time.Hour * 24,
			},
			want: 24,
		},
		{
			name: "Days of Month",
			args: args{
				from: &Timestamp{
					t: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				to: &Timestamp{
					t: time.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC),
				},
				d: time.Hour * 24,
			},
			want: 31,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ToRange(tt.args.from, tt.args.to, tt.args.d)
			for _, v := range r {
				fmt.Println(v.String())
			}

			//if got := ToRange(tt.args.from, tt.args.to, tt.args.interval, tt.args.d); !reflect.DeepEqual(len(got), tt.want) {
			//	t.Errorf("ToRange() = %v, want %v", len(got), tt.want)
			//}
		})
	}
}

func TestAddMonth(t *testing.T) {
	type args struct {
		ts *Timestamp
		m  time.Month
	}
	tests := []struct {
		name string
		args args
		want Timestamp
	}{
		{
			name: "test add 1 month from October 31",
			args: args{
				ts: &Timestamp{
					t: time.Date(2024, time.October, 31, 1, 2, 3, 0, time.UTC),
				},
				m: 1,
			},
			want: Timestamp{
				t: time.Date(2024, time.November, 30, 1, 2, 3, 0, time.UTC),
			},
		},
		{
			name: "test add 1 month from December 31, 2024",
			args: args{
				ts: &Timestamp{
					t: time.Date(2024, time.December, 31, 1, 2, 3, 0, time.UTC),
				},
				m: 1,
			},
			want: Timestamp{
				t: time.Date(2025, time.January, 31, 1, 2, 3, 0, time.UTC),
			},
		},
		{
			name: "test add 1 month from February 29",
			args: args{
				ts: &Timestamp{
					t: time.Date(2024, time.February, 29, 1, 2, 3, 0, time.UTC),
				},
				m: 1,
			},
			want: Timestamp{
				t: time.Date(2024, time.March, 29, 1, 2, 3, 0, time.UTC),
			},
		},
		{
			name: "test add 1 month from January 29 on a leap year",
			args: args{
				ts: &Timestamp{
					t: time.Date(2024, time.January, 29, 1, 2, 3, 0, time.UTC),
				},
				m: 1,
			},
			want: Timestamp{
				t: time.Date(2024, time.February, 29, 1, 2, 3, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddMonth(tt.args.ts, int(tt.args.m)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMonthToPeriod(t *testing.T) {
	type args struct {
		ts *Timestamp
	}
	tests := []struct {
		name string
		args args
		want Period
	}{
		{
			name: "test January 2024",
			args: args{
				ts: &Timestamp{
					t: time.Date(2024, time.January, 15, 1, 2, 3, 0, time.UTC),
				},
			},
			want: Period{
				Start: &Timestamp{
					t: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				End: &Timestamp{
					t: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "test February 2024",
			args: args{
				ts: &Timestamp{
					t: time.Date(2024, time.February, 15, 1, 2, 3, 0, time.UTC),
				},
			},
			want: Period{
				Start: &Timestamp{
					t: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
				},
				End: &Timestamp{
					t: time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "test February 2023",
			args: args{
				ts: &Timestamp{
					t: time.Date(2023, time.February, 15, 1, 2, 3, 0, time.UTC),
				},
			},
			want: Period{
				Start: &Timestamp{
					t: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
				},
				End: &Timestamp{
					t: time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MonthToPeriod(tt.args.ts.Year(), tt.args.ts.t.Month()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MonthToPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestISO8601ToDuration(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  time.Duration
		expectErr bool
	}{
		// Valid Test Cases
		{"Complete ISO8601 Duration", "P1Y2M3W4DT5H6M7S", time.Duration((1 * 365 * 24 * time.Hour) + (2 * 30 * 24 * time.Hour) + (3 * 7 * 24 * time.Hour) + (4 * 24 * time.Hour) + (5 * time.Hour) + (6 * time.Minute) + (7 * time.Second)), false},
		{"Years Only", "P2Y", time.Duration(2 * 365 * 24 * time.Hour), false},
		{"Months Only", "P3M", time.Duration(3 * 30 * 24 * time.Hour), false},
		{"Weeks Only", "P4W", time.Duration(4 * 7 * 24 * time.Hour), false},
		{"Days Only", "P5D", time.Duration(5 * 24 * time.Hour), false},
		{"Hours Only", "PT6H", time.Duration(6 * time.Hour), false},
		{"Minutes Only", "PT7M", time.Duration(7 * time.Minute), false},
		{"Seconds Only", "PT30S", time.Duration(30 * time.Second), false},
		{"Decimal Seconds Stripped", "PT5.9S", time.Duration(5 * time.Second), false},
		{"Mixed Time Only", "PT1H30M45S", time.Duration(1*time.Hour + 30*time.Minute + 45*time.Second), false},
		{"Mixed Date Only", "P1Y2M10D", time.Duration(1*365*24*time.Hour + 2*30*24*time.Hour + 10*24*time.Hour), false},
		{"Empty Time Designator", "P1Y2M3W4D", time.Duration(1*365*24*time.Hour + 2*30*24*time.Hour + 3*7*24*time.Hour + 4*24*time.Hour), false},
		{"Only Time Designator", "PT", time.Duration(0), true},
		{"Only Time Designator", "PT40M3S", time.Duration(40*time.Minute + 3*time.Second), false},

		// Invalid Test Cases
		{"Empty String", "", time.Duration(0), true},
		{"Invalid Format", "1Y2M3W4DT5H6M7S", time.Duration(0), true},
		{"Missing Designator", "P1Y2M3DT", time.Duration(0), true},
		{"Negative Input", "PT-5S", time.Duration(0), true},
		{"Non-Numeric Values", "P1Y2M3W4DT5H6MXS", time.Duration(0), true},
		{"Only Non-ISO8601 Format String", "AnythingNotISO", time.Duration(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dur, err := ParseISO8601Duration(tt.input)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error, got: %v", err)
				}
				if dur != tt.expected {
					t.Errorf("expected duration %v, got %v", tt.expected, dur)
				}
			}
		})
	}
}
