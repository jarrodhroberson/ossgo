package seq

import (
	"cmp"
	"iter"
	"slices"
	"sync/atomic"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

func Empty[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {
		return
	}
}

func Empty2[K comparable, V any]() iter.Seq2[K,V] {
	return func(yield func(K, V) bool) {
		return
	}
}

func Collect2[K comparable, V any](it iter.Seq2[K, V]) map[K]V {
	m := make(map[K]V)
	for k, v := range it {
		m[k] = v
	}
	return m
}

func Filter[T any](it iter.Seq[T], predicate func(i T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range it {
			if predicate(i) {
				if !yield(i) {
					return
				}
			}
		}
	}
}

func Filter2[K any, V any](it iter.Seq2[K, V], predicate func(k K, v V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range it {
			if predicate(k, v) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// WrapWithCounter takes an iter.Seq and returns a *CountingSeq.
// The *CountingSeq will count the number of items that have been yielded
// when the Seq function is called.
// This is useful for debugging and testing.
func WrapWithCounter[T any](it iter.Seq[T]) *CountingSeq[T] {
	var counter atomic.Int64
	return &CountingSeq[T]{
		Seq: func(yield func(T) bool) {
			for i := range it {
				if !yield(i) {
					return
				}
				counter.Add(1)
			}
		},
		counter: &counter,
	}
}

// IntRange returns an iter.Seq that yields a sequence of integers from start to end inclusive.
//
//	for i := range IntRange(1, 5) {
//		fmt.Println(i) // 1, 2, 3, 4, 5
//	}
func IntRange[N Integer](start N, end N) iter.Seq[N] {
	return func(yield func(N) bool) {
		for i := start; i <= end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// ToSeq takes a slice of type T and returns an iter.Seq[T] that yields each item in the slice.
//
//	for i := range ToSeq([]int{1, 2, 3}) {
//		fmt.Println(i) // 1, 2, 3
//	}
func ToSeq[T any](seq ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range seq {
			if !yield(seq[i]) {
				return
			}
		}
	}
}

// PassThruFunc passes thru the value unchanged
// this is just a convenience function for times when you do not want to transform the key or value in Map2
// so you do not have to write an inline function and clutter up the code more than it needs to be.
func PassThruFunc[K any](k K) K {
	return k
}

// Map also known as transform function, list comprehension, visitor pattern whatever.
// this takes an iter.Seq and a function to apply to each item in the sequence
// returning a new iter.Seq with the results.
func Map[T any, R any](it iter.Seq[T], mapFunc func(t T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for i := range it {
			if !yield(mapFunc(i)) {
				return
			}
		}
	}
}

// Map2 this does the same but allows you to transform the key and/or the value.
func Map2[K any, V any, KR any, VR any](it iter.Seq2[K, V], keyFunc func(k K) KR, valFunc func(v V) VR) iter.Seq2[KR, VR] {
	return func(yield func(KR, VR) bool) {
		for k, v := range it {
			if !yield(keyFunc(k), valFunc(v)) {
				return
			}
		}
	}
}

// ToSeq2 converts an iter.Seq[V] to an iter.Seq2[K, V] by applying a key function to each value.
// This is useful when you have a sequence of values and you want to associate a key with each value.
// The key function is applied to each value in the sequence to generate the key.
// The original value is preserved as the value in the resulting iter.Seq2.
func ToSeq2[K any, V any](is iter.Seq[V], keyFunc func(v V) K) iter.Seq2[K, V] {
	return iter.Seq2[K, V](func(yield func(K, V) bool) {
		for v := range is {
			k := keyFunc(v)
			if !yield(k, v) {
				return
			}
		}
	})
}

// SkipAndLimit combines SkipFirstN and FirstN in a way that allows accessing sub sequences without any overhead.
// it first SkipFirstN items using skip and then returns the next FirstN or less items in limit.
// using this instead of inlining it shows more intent because of the semantic of the name.
// if skip is greater than the number of items in the sequence then an empty sequence is returned.
// if limit is greater than the number of items remaining in the sequence then the sequence will
// contain less than limit number of items in it.
func SkipAndLimit[V any](it iter.Seq[V], skip int, limit int) iter.Seq[V] {
	return FirstN[V](SkipFirstN[V](it, skip), limit)
}

// FirstN takes an iter.Seq[int] and returns a new iter.Seq[int] that
// yields only the first 'limit' or less items without creating intermediate slices.
// if there are less than limit items in the sequence then the sequence ends
// when the sequence end is reached.
func FirstN[T any](it iter.Seq[T], limit int) iter.Seq[T] {
	return iter.Seq[T](func(yield func(T) bool) {
		count := 0
		for item := range it {
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

// First returns the first item in the sequence that matches the predicate.
// if no item matches the predicate then the sequence is empty.
func FindFirst[T any](it iter.Seq[T], predicate func(i T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range it {
			if predicate(i) {
				if !yield(i) {
					return
				}
				return
			}
		}
	}
}

// SkipFirstN skips the first N items in the sequence and then iterates over the rest of them normally.
// if skip is greater than the number of items in the sequence then an empty sequence is returned.
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
// All but the last iter.Seq chunk will have size n.
// Chunk panics if n is less than 1.
func Chunk[T any](sq iter.Seq[T], size int) iter.Seq[iter.Seq[T]] {
	if size < 0 {
		panic(errs.MinSizeExceededError.New("size %d must be >= 0", size))
	}

	return func(yield func(s iter.Seq[T]) bool) {
		next, stop := iter.Pull[T](sq)
		defer stop()
		endOfSeq := false
		for !endOfSeq {
			// get the first item for the chunk
			v, ok := next()
			// there are no more items !ok then exit loop
			// this prevents returning an extra empty iter.Seq at end of Seq
			if !ok {
				break
			}
			// create the next sequence chunk
			iterSeqChunk := func(yield func(T) bool) {
				i := 0
				for ; i < size; i++ {
					if ok {
						if !ok {
							// end of original sequence
							// this sequence may be <= size
							endOfSeq = true
							break
						}

						if !yield(v) {
							return
						}
						v, ok = next()
					}
				}
			}
			if !yield(iterSeqChunk) {
				return
			}
		}
	}
}

// Chunk2 returns an iterator over consecutive sub-slices of up to n elements of s.
// All but the last iter.Seq chunk will have size n.
// Chunk2 panics if n is less than 1.
func Chunk2[K any, V any](sq iter.Seq2[K, V], size int) iter.Seq[iter.Seq2[K, V]] {
	if size < 0 {
		panic(errs.MinSizeExceededError.New("size %d must be >= 0", size))
	}

	return func(yield func(s iter.Seq2[K, V]) bool) {
		next, stop := iter.Pull2[K, V](sq)
		defer stop()
		endOfSeq := false
		for !endOfSeq {
			// get the first item for the chunk
			k, v, ok := next()
			// there are no more items !ok then exit loop
			// this prevents returning an extra empty iter.Seq at end of Seq
			if !ok {
				break
			}
			// create the next sequence chunk
			iterSeqChunk := func(yield func(K, V) bool) {
				i := 0
				for ; i < size; i++ {
					if ok {
						if !ok {
							// end of original sequence
							// this sequence may be <= size
							endOfSeq = true
							break
						}

						if !yield(k, v) {
							return
						}
						k, v, ok = next()
					}
				}
			}
			if !yield(iterSeqChunk) {
				return
			}
		}
	}
}

// Min returns the minimum value in the sequence.
// It uses the cmp.Ordered constraint to compare values.
// It assumes the sequence is not empty.
// If the sequence is empty it will panic.
func Min[T cmp.Ordered](it iter.Seq[T]) T {
	minV := slices.Collect(FirstN[T](it, 1))[0]
	for v := range it {
		minV = min(minV, v)
	}
	return minV
}

// First returns the first item in the sequence.
// If the sequence is empty, it returns a errs.NotFoundError.
func First[T any](it iter.Seq[T]) (T, bool) {
	next, stop := iter.Pull[T](it)
	defer stop()
	return next()
}
