package seq

import (
	"iter"
	"sync/atomic"
)

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type Decimal interface {
	float32 | float64
}

// Number is a union of Integer and Decimal.
type Number interface {
	Integer | Decimal
}

// CountingSeq is a wraps an iter.Seq and counts the number of elements that have been
// iterated over.
type CountingSeq[T any] struct {
	Seq     iter.Seq[T]
	counter *atomic.Int64
}

func (c CountingSeq[T]) Count() int64 {
	return c.counter.Load()
}

type MemoizedSeq[T any] struct {
	Seq   iter.Seq[T]
	reset func()
}

func (m *MemoizedSeq[T]) Reset() {
	m.reset()
}

// ToMemoizingSeq wraps an iter.Seq to memoize its items and allows it to be re-iterated after a reset() call.
func ToMemoizingSeq[T any](seq iter.Seq[T]) MemoizedSeq[T] {
	var items []T
	var memoized bool

	memoizedSeq := func(yield func(item T) bool) {
		if memoized {
			for _, item := range items {
				if !yield(item) {
					return
				}
			}
		} else {
			seq(func(item T) bool {
				items = append(items, item)
				return yield(item)
			})
			memoized = true
		}
	}

	reset := func() {
		memoized = false
	}

	return MemoizedSeq[T]{
		Seq:   memoizedSeq,
		reset: reset,
	}
}

// FlattenSeq takes an iter.Seq of batches (iter.Seq[T]) and flat maps all the batches
// into a single iter.Seq.
func FlattenSeq[T any](iterSeqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(t T) bool) {
		for is := range iterSeqs {
			for i := range iterSeqs[is] {
				if !yield(i) {
					return
				}
			}
		}
	}
}
