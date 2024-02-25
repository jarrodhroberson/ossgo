package timestamp

import (
	"reflect"
	"testing"
	"time"
)

func TestToRange(t *testing.T) {
	type args struct {
		from Timestamp
		to   Timestamp
		d    time.Duration
	}
	tests := []struct {
		name string
		args args
		want []Timestamp
	}{
		{
			name: "test range",
			args: args{
				from: MustParse(time.RFC3339Nano, "2009-06-09T03:55:33Z"),
				to:   MustParse(time.RFC3339Nano, "2010-01-01T00:00:00Z"),
				d:    (time.Hour * 24) * 7,
			},
			want: []Timestamp{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToRange(tt.args.from, tt.args.to, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
