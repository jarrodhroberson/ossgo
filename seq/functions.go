package seq

import (
	"iter"

	errs "github.com/jarrodhroberson/ossgo/errors"
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

func Map[T any, R any](it iter.Seq[T], mapFunc func(t T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for i := range it {
			if !yield(mapFunc(i)) {
				return
			}
		}
	}
}

func SeqToSeq2[K any, V any](is iter.Seq[V], keyFunc func(v V) K) iter.Seq2[K, V] {
	return iter.Seq2[K, V](func(yield func(K, V) bool) {
		for v := range is {
			k := keyFunc(v)
			if !yield(k, v) {
				return
			}
		}
	})
}

func SkipLimit[V any](it iter.Seq[V], skip int, limit int) iter.Seq[V] {
	return FirstN[V](SkipFirstN[V](it, skip), limit)
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
				for ; i < size-1; i++ {
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
				for ; i < size-1; i++ {
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
