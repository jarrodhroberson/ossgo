package seq

import (
	"fmt"
	"iter"
)

func IntRange(start int, end int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := start; i <= end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// FirstN takes an iter.Seq[int] and returns a new iter.Seq[int] that
// yields only the first 'limit' items without creating intermediate slices.
func FirstN[T any](original iter.Seq[T], limit int) iter.Seq[T] {
	return iter.Seq[T](func(yield func(T) bool) {
		count := 0
		for item := range original {
			if count < limit {
				if !yield(item) {
					return
				}
				count++
			} else {
				return
			}
		}
	})
}

func SkipFirstN[T any](seq iter.Seq[T], skip int) iter.Seq[T] {
	return iter.Seq[T](func(yield func(T) bool) {
		next, stop := iter.Pull[T](seq)
		defer stop()

		for i := 0; i <= skip; i++ {
			_, ok := next()
			if !ok {
				break
			}
		}
		for {
			v, ok := next()
			if !ok {
				break
			}
			if !yield(v) {
				return
			}
		}
	})
}

// Chunk returns an iterator over consecutive sub-slices of up to n elements of s.
// All but the last sub-slice will have size n.
// All sub-slices are clipped to have no capacity beyond the length.
// If s is empty, the sequence is empty: there is no empty slice in the sequence.
// Chunk panics if n is less than 1.
func Chunk[E any](s1 iter.Seq[E], n int) iter.Seq[iter.Seq[E]] {
	start := 0
	return func(yield func(s iter.Seq[E]) bool) {
		for i := start; i < n; i++ {
			fmt.Printf("start: %d\n", i)
			s2 := FirstN(SkipFirstN(s1, n), n)
			if !yield(s2) {
				return
			}
		}
		start = start + n
	}
}
