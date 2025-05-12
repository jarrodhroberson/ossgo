package seq

import (
    "bufio"
    "cmp"
    "io"
    "iter"
    "maps"
    "slices"
    "sync"
    "sync/atomic"

    "github.com/jarrodhroberson/destruct/destruct"
    "github.com/rs/zerolog/log"

    errs "github.com/jarrodhroberson/ossgo/errors"
    "github.com/jarrodhroberson/ossgo/functions/must"
)

// Empty returns an iter.Seq[T] that does not yield any items.
//
// Example usage:
//
//	for i := range seq.Empty[int]() {
//		fmt.Println(i) // This will not print anything as the sequence is empty
//	}
//
// This can be useful when you need to provide an empty sequence to a function
// that expects an iter.Seq[T], without requiring any special case handling for the caller.
func Empty[T any]() iter.Seq[T] {
    return func(yield func(T) bool) {
        return
    }
}

// Empty2 returns an iter.Seq2[K, V] that does not yield any items.
//
// Example usage:
//
//	for k, v := range seq.Empty2[int, string]() {
//		fmt.Println(k, v) // This will not print anything as the sequence is empty
//	}
//
// This can be useful when you need to provide an empty sequence to a function
// that expects an iter.Seq2[K, V], without requiring any special case handling for the caller.
func Empty2[K comparable, V any]() iter.Seq2[K, V] {
    return func(yield func(K, V) bool) {
        return
    }
}

// Collect2 collects all key-value pairs from an iter.Seq2[K, V] into a map[K]V.
//
// Example usage:
//
//	seq := seq.ToSeq2(seq.ToSeq(1, 2, 3), func(v int) string { return fmt.Sprintf("Key%d", v) })
//	result := seq.Collect2(seq)
//
//	fmt.Println(result) // Output: map[Key1:1 Key2:2 Key3:3]
func Collect2[K comparable, V any](it iter.Seq2[K, V]) map[K]V {
    m := make(map[K]V)
    for k, v := range it {
        m[k] = v
    }
    return m
}

// Filter filters elements from an iter.Seq[T] that match the given predicate function.
//
// Parameters:
//   - it: The input sequence of type T.
//   - predicate: A function that takes an element of type T and returns a boolean indicating
//     whether the element should be included in the resulting sequence.
//
// Returns:
//   - An iter.Seq[T] containing only those elements of the input sequence that satisfy the predicate.
//
// Example usage:
//
//	seq := seq.ToSeq(1, 2, 3, 4, 5)
//	filtered := seq.Filter(seq, func(i int) bool { return i%2 == 0 })
//
//	for v := range filtered {
//		fmt.Println(v) // Output: 2, 4
//	}
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

// Filter2 filters key-value pairs from an iter.Seq2[K, V] that satisfy the given predicate function.
//
// Parameters:
//   - it: The input sequence of key-value pairs.
//   - predicate: A function that takes a key and a value, returning true if the pair
//     should be included in the resulting sequence and false otherwise.
//
// Returns:
//   - An iter.Seq2[K, V] containing only those key-value pairs of the input sequence
//     that satisfy the predicate.
//
// Example usage:
//
//	pairs := seq.ToSeq2(seq.ToSeq(1, 2, 3), func(v int) string { return fmt.Sprintf("Key%d", v) })
//	filtered := seq.Filter2(pairs, func(k string, v int) bool { return v%2 == 1 })
//
//	for k, v := range filtered {
//		fmt.Println(k, v) // Output: Key1 1, Key3 3
//	}
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
// If the sequence is empty, it returns nil, false.
func First[T any](it iter.Seq[T]) (T, bool) {
    next, stop := iter.Pull[T](it)
    defer stop()
    return next()
}

// RuneSeq creates a sequence of runes from the given string.
//
// The sequence iterates over all runes in the input string and passes each rune
// to the provided yielding function until all runes are processed, or the yielding
// function returns false to stop the iteration.
//
// Parameters:
//   - s: The input string to generate the rune sequence from.
//
// Returns:
//   - iter.Seq[rune]: A sequence that yields runes from the input string.
//
// Example:
//
//	  seq := RuneSeq("hello")
//	  seq(func(r rune) bool {
//		   fmt.Printf("Rune: %c\n", r)
//		   return true // Continue iteration
//	  })
func RuneSeq(s string) iter.Seq[rune] {
    return func(yield func(rune) bool) {
        for _, r := range s {
            if !yield(r) {
                return
            }
        }
    }
}

// RuneSeq2 creates a sequence of runes from the given string,
// where each yielded element is a key-value pair consisting of the index
// and the rune at that index in the string.
//
// The sequence will iterate over the string, yielding the index and
// corresponding rune until all runes in the string have been processed,
// or until the yielding function returns false to stop iteration.
//
// Parameters:
//   - s: The input string to generate the rune sequence from.
//
// Returns:
//   - iter.Seq2[int,rune]: A sequence that yields index-rune pairs from the string.
//
// Example:
//
//	  seq := RuneSeq2("hello")
//	  seq(func(idx int, r rune) bool {
//		   fmt.Printf("Index: %d, Rune: %c\n", idx, r)
//		   return true // Continue iteration
//	  })
func RuneSeq2(s string) iter.Seq2[int, rune] {
    return func(yield func(int, rune) bool) {
        for idx, r := range s {
            if !yield(idx, r) {
                return
            }
        }
    }
}

// OrderedIterSeq creates an ordered sequence from the given input sequence.
// It reads all elements from the input sequence into memory, sorts them using
// their natural order (defined by the cmp.Ordered constraint), and produces
// an iterator that yields the sorted elements.
//
// The input sequence is consumed fully before yielding the sorted sequence,
// so this function is suitable for sequences that can fit into memory.
//
// Parameters:
//   - in: An input sequence of elements implementing the cmp.Ordered constraint.
//
// Returns:
//   - iter.Seq[T]: An iterator that yields the elements of the input sequence
//     in ascending order.
//
// Usage:
//
//	  seq := OrderedIterSeq(SomeSeq)
//	  seq(func(v int) bool {
//		   fmt.Println(v)
//		   return true // Continue iteration
//	  })
//
// Notes:
//   - The ordering is determined by slices.Sort, which uses the natural order of the items.
//   - The function utilizes a buffered channel of size 256 to pass sorted items to the output sequence.
//   - Synchronization is ensured using a WaitGroup and a done channel to prevent race conditions.
func OrderedIterSeq[T cmp.Ordered](in iter.Seq[T]) iter.Seq[T] {
    orderChan := make(chan T, 256)
    done := make(chan struct{})
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        defer close(orderChan)
        defer wg.Done()

        var items []T
        for item := range in {
            items = append(items, item)
        }

        slices.Sort(items)

        for item := range ToSeq(items...) {
            orderChan <- item
        }

        close(done)
    }()

    return iter.Seq[T](func(yield func(T) bool) {
        defer wg.Wait()

        for {
            select {
            case item, ok := <-orderChan:
                if !ok {
                    return
                }
                if !yield(item) {
                    return
                }
            case <-done:
                return
            }
        }
    })
}

// Reduce reduces a sequence to a single value by repeatedly applying a function
// to an accumulator and each element of the sequence.
//
// Parameters:
//   - s: An input sequence of elements.
//   - initialValue: The initial value of the accumulator.
//   - f: A function that combines the accumulator with an element from the sequence.
//
// Returns:
//   - The final accumulated value after processing the entire sequence.
//
// Example:
//
//	  seq := iter.Seq[int](func(yield func(int) bool) {
//		   yield(1)
//		   yield(2)
//		   yield(3)
//	  })
//
//	  sum := Reduce(seq, 0, func(acc, val int) int {
//		   return acc + val
//	  })
//	  fmt.Println(sum) // Output: 6
//
// Notes:
//   - The sequence `s` will be consumed entirely by this function.
//   - The function `f` must be capable of reducing any two elements into the accumulator.
func Reduce[T any, A any](s iter.Seq[T], initialValue A, f func(A, T) A) A {
    acc := initialValue
    for i := range s {
        acc = f(acc, i)
    }
    return acc
}

// Reduce2 reduces a sequence of key-value pairs to a single value by repeatedly applying a function
// to an accumulator and each key-value pair from the sequence.
//
// Parameters:
//   - s: An input sequence of key-value pairs.
//   - initialValue: The initial value of the accumulator.
//   - f: A function that combines the accumulator with a key and a value from the sequence.
//
// Returns:
//   - The final accumulated value after processing the entire sequence.
//
// Example:
//
//	  seq := iter.Seq2[int, string](func(yield func(int, string) bool) {
//		   yield(1, "a")
//		   yield(2, "b")
//		   yield(3, "c")
//	  })
//
//	  result := Reduce2(seq, "", func(acc string, k int, v string) string {
//		   return acc + fmt.Sprintf("%d:%s ", k, v)
//	  })
//	  fmt.Println(result) // Output: "1:a 2:b 3:c "
//
// Notes:
//   - The sequence `s` will be consumed entirely by this function.
//   - The function `f` must be capable of reducing the key-value pairs into the accumulator.
func Reduce2[K any, V any, A any](s iter.Seq2[K, V], initialValue A, f func(A, K, V) A) A {
    acc := initialValue
    for k, v := range s {
        acc = f(acc, k, v)
    }
    return acc
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

// Sum calculates the sum of a sequence of numbers.
//
// Parameters:
//   - s: An input sequence of numeric elements.
//
// Returns:
//   - The sum of all elements in the sequence.
//
// Example:
//
//	  seq := iter.Seq[int](func(yield func(int) bool) {
//		   yield(1)
//		   yield(2)
//		   yield(3)
//	  })
//
//	  total := Sum(seq)
//	  fmt.Println(total) // Output: 6
func Sum[T Number](s iter.Seq[T]) T {
    return Reduce[T, T](s, 0, func(a T, t T) T {
        return a + t
    })
}

// Unique filters a sequence to include only unique elements based on their hash identity.
//
// It ensures that only one instance of each unique element is yielded, with uniqueness
// determined by a hash of the element.
//
// Parameters:
//   - s: An input sequence of elements.
//
// Returns:
//   - A new sequence containing only unique elements from the input sequence.
//
// Notes:
//   - The uniqueness of elements is determined using their hash identity via MustHashIdentity.
//   - If elements cannot be hashed correctly, this function might panic.
//
// Example:
//
//	  seq := iter.Seq[int](func(yield func(int) bool) {
//			 yield(1)
//			 yield(2)
//			 yield(2)
//			 yield(3)
//	  })
//
//	  uniqueSeq := Unique(seq)
//	  for v := range uniqueSeq {
//		   fmt.Println(v) // Output: 1, 2, 3
//	  }
func Unique[T any](s iter.Seq[T]) iter.Seq[T] {
    return func(yield func(T) bool) {
        seen := make(map[string]struct{}) // Use a map to track seen elements.

        for v := range s {
            key := destruct.MustHashIdentity(v)
            if _, ok := seen[key]; ok {
                continue
            } else {
                if !yield(v) {
                    return
                }
                seen[key] = struct{}{}
            }

        }
    }

}

// Count counts the number of elements in a sequence.
//
// Parameters:
//   - s: An input sequence of elements.
//
// Returns:
//   - The total count of elements in the sequence.
//
// Example:
//
//	  seq := iter.Seq[int](func(yield func(int) bool) {
//		   yield(1)
//		   yield(2)
//		   yield(3)
//	  })
//
//	  count := Count(seq)
//	  fmt.Println(count) // Output: 3
//
// Notes:
//   - The sequence `s` will be completely consumed by this function.
func Count[T any](s iter.Seq[T]) int64 {
    return Reduce[T, int64](s, 0, func(acc int64, _ T) int64 {
        return acc + 1
    })
}

// Count2 counts the number of key-value pairs in a sequence of pairs.
//
// Parameters:
//   - s: An input sequence of key-value pairs.
//
// Returns:
//   - The total count of key-value pairs in the sequence.
//
// Example:
//
//	  seq := iter.Seq2[int, string](func(yield func(int, string) bool) {
//		   yield(1, "a")
//		   yield(2, "b")
//		   yield(3, "c")
//	  })
//
//	  count := Count2(seq)
//	  fmt.Println(count) // Output: 3
//
// Notes:
//   - The sequence `s` will be completely consumed by this function.
func Count2[K any, V any](s iter.Seq2[K, V]) int64 {
    return Reduce2[K, V, int64](s, int64(0), func(acc int64, k K, v V) int64 {
        return acc + 1
    })
}

// Seq2ToMap converts an iter.Seq2[string, any] to a map[string]any.
//
// Parameters:
//   - seq: An input sequence of key-value pairs.
//
// Returns:
//   - A map[string]any containing the key-value pairs from the input sequence.
//
// this just delegates to maps.Collect() because I keep forgetting it exists
// Deprecated: Use maps.Collect instead for more streamlined functionality.
func Seq2ToMap[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
    return maps.Collect(seq)
}

// UnzipMap splits a map into two slices: one containing its keys and the other its values.
//
// Parameters:
//   - m: The input map of type map[K]V.
//
// Returns:
//   - A slice of keys ([]K) and a slice of values ([]V) extracted from the input map,
//     sorted by the keys to ensure consistency in ordering.
//
// Example:
//
//	  m := map[string]int{
//		   "a": 1,
//		   "b": 2,
//		   "c": 3,
//	  }
//
//	  keys, values := UnzipMap(m)
//	  fmt.Println(keys)   // Output: [a b c]
//	  fmt.Println(values) // Output: [1 2 3]
func UnzipMap[K comparable, V any](m map[K]V) ([]K, []V) {
    keys := make([]K, 0, len(m))
    values := make([]V, 0, len(m))

    for k := range maps.Keys(m) {
        keys = append(keys, k)
        values = append(values, m[k])
    }

    return keys, values
}

// GroupBy groups elements of a sequence by a key function. It collects elements
// that share the same key into separate sequences.
//
// Parameters:
//   - s: The input sequence of type iter.Seq[V].
//   - keyFunc: A function that extracts the key (of type K) for each element.
//   - groupByFunc: A function that determines if an element belongs to a specific group.
//
// Returns:
//   - An iter.Seq2[K, iter.Seq[V]], where each key (K) is associated with a sequence
//     of elements (iter.Seq[V]) that belong to its group.
//
// Example:
//
//	  seq := iter.Seq[int](func(yield func(int) bool) {
//		   yield(1)
//		   yield(2)
//		   yield(3)
//		   yield(4)
//		   yield(5)
//	  })
//	  grouped := GroupBy(seq, func(v int) int { return v % 2 }, func(k, v int) bool { return v % 2 == k })
//	  for key, group := range grouped {
//		   fmt.Println("Key:", key)
//		   for item := range group {
//			   fmt.Println("  Item:", item)
//		   }
//	  }
//	  // Output:
//	  // Key: 0
//	  //   Item: 2
//	  //   Item: 4
//	  // Key: 1
//	  //   Item: 1
//	  //   Item: 3
//	  //   Item: 5
func GroupBy[K comparable, V any](s iter.Seq[V], keyFunc func(V) K, groupByFunc func(K, V) bool) iter.Seq2[K, iter.Seq[V]] {
    return func(yield func(K, iter.Seq[V]) bool) {
        groupKeys := Unique[K](Map[V, K](s, keyFunc))
        for groupKey := range groupKeys {
            group := Filter[V](s, func(v V) bool {
                return groupByFunc(groupKey, v)
            })
            iter2 := func(yield func(V) bool) {
                for item := range group {
                    if !yield(item) {
                        return
                    }
                }
            }
            if !yield(groupKey, iter2) {
                return
            }
        }
    }
}

//
// JsonMarshalCompact serializes the given sequence into a compact JSON array
// and writes it to the provided writer.
//
// Parameters:
//   - w: An `io.Writer` where the serialized JSON array is written.
//   - seq: A sequence `iter.Seq[T]` of elements to be serialized.
//
// Returns:
//   - An error if any issue occurs during marshaling or writing.
//
// Notes:
//   - Uses a buffered writer for efficient writing, ensuring any pending 
//     data is flushed at the end.
//   - Serializes the sequence items one by one, maintaining a compact JSON format  
//     without extra whitespace.
//
// Example:
//
//   seq := iter.Seq[int](func(yield func(int) bool) {
//       yield(1)
//       yield(2)
//       yield(3)
//   })
//
//   var buf bytes.Buffer
//   err := JsonMarshalCompact(&buf, seq)
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(buf.String()) // Output: [1,2,3]
//
//   If the writer fails at any point, marshaling stops and an error is returned.
func JsonMarshalCompact[T any](w io.Writer, seq iter.Seq[T]) error {
    // no need to check if it is already a buffered writer, it does that already
    // default size is 4096 bytes
    w = bufio.NewWriter(w)
    defer func() {
        if err := w.(*bufio.Writer).Flush(); err != nil {
            err = errs.MarshalError.New("failed to flush buffered writer: %w", err)
            log.Warn().Err(err).Msg(err.Error())
        }
    }()

    first := true
    next, stop := iter.Pull(seq)
    defer stop()
    for {
        if t, ok := next(); ok {
            if first {
                if _, err := w.Write([]byte{'['}); err != nil {
                    return errs.MarshalError.New("failed to write opening bracket: %w", err)
                }
                first = false
            } else {
                if _, err := w.Write([]byte{','}); err != nil {
                    err = errs.MarshalError.New("failed to write comma: %w", err)
                    return err // Stop iteration on write error
                }
                _, err := w.Write(must.MarshalJson(t))
                if err != nil {
                    return err
                }
            }
        }
        if _, err := w.Write([]byte{']'}); err != nil {
            return errs.MarshalError.New("failed to write closing bracket: %w", err)
        }
    }
}
