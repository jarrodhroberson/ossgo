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

// PassThruFunc passes thru the value unchanged
// this is just a convience function for times when you do not want to transform the key or value in Map2
// so you do not have to write an inline function and clutter up the code more than it needs to be.
func PassThruFunc[K any](k K) K {
	return k
}

// Map also know as transform function, list comprehension, visitor pattern whatever.
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

// SkipAndLimit combines SkipFirstN and FirstN in a way that allows accessing sub sequences without any overhead.
// it first SkipFirstN items using skip and then returns the next FirstN items in limit.
// using this instead of inlining it shows more intent because of the semantic of the name.
func SkipAndLimit[V any](it iter.Seq[V], skip int, limit int) iter.Seq[V] {
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

// SkipFirstN skips the first N items in the sequence and then iterates over the rest of them normally.
// SkipAndLimit combines this and FirstN to provide arbitrary sub sequencing zero cost abstraction.
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
