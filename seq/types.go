package seq

import (
	"iter"
	"slices"
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

type MemoizeSeq[T any] struct {
	delegate []T
}

func (m *MemoizeSeq[T]) Seq() iter.Seq[T] {
	return slices.Values(m.delegate)
}

func (m *MemoizeSeq[T]) Seq2() iter.Seq2[int,T] {
	return slices.All(m.delegate)
}

func (m *MemoizeSeq[T]) Len() int {
	return len(m.delegate)
}