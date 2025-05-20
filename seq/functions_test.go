package seq

import (
	"iter"
	"testing"
)

func TestCount(t *testing.T) {
	type args[T any] struct {
		s iter.Seq[int]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want int64
	}
	tests := []testCase[int]{
		{
			name: "count_0_to_9",
			args: args[int]{
				s: IntRange(0, 9),
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Count(tt.args.s); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}
