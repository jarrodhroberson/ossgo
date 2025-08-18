package timestamp

import (
	"reflect"
	"testing"
	"time"
)

func TestDayToPeriod(t *testing.T) {
	type args struct {
		ts *Timestamp
	}
	tests := []struct {
		name string
		args args
		want *Period
	}{
		{
			name: "test_with_valid_timestamp",
			args: args{
				ts: &Timestamp{
					t: time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
				},
			},
			want: &Period{
				start: &Timestamp{t: time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC)},
				end: &Timestamp{
					t: time.Date(2023, 10, 6, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DayToPeriod(tt.args.ts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayToPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}
