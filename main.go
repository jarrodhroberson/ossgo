package main

import (
	"fmt"

	"github.com/jarrodhroberson/ossgo/seq"
)

func main() {
	fmt.Println("ossgo")
	testSeq()
}

func testSeq() {
	ints := seq.IntRange[int64](100, 199)
	//next, _ := iter.Pull(seq.Sum(ints))
	//sum, _ := next()
	//fmt.Println(sum)

	count := 0
	for cis := range seq.Chunk(ints, 10) {
		fmt.Printf("%d:", count)
		for i := range cis {
			fmt.Printf("%02d,", i)
		}
		fmt.Println()
		count++
	}

	//next, stop := iter.Pull(func(start int) iter.Seq[int] {
	//	index := start
	//	return iter.Seq[int](func(yield func(int) bool) {
	//		for {
	//			if !yield(index) {
	//				return
	//			}
	//			index++
	//		}
	//	})
	//}(0))
	//
	//intMap := seq.SeqToSeq2[int, int](ints, func(v int) int {
	//	k, ok := next()
	//	if !ok {
	//		stop()
	//		return -1
	//	}
	//	return k
	//})
	//
	//count := 0
	//for cis2 := range seq.Chunk2(intMap, 10) {
	//	fmt.Printf("%d:", count)
	//	for k, v := range cis2 {
	//		fmt.Printf("[%02d,%02d],", k, v)
	//	}
	//	fmt.Println()
	//	count++
	//}
}
