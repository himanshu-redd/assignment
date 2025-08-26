package main

import (
	"fmt"
	"sort"
)

type Pair struct {
	num  int
	freq int
}

func topKFrequent(nums []int, k int) []int {
	freq := make(map[int]int)
	for _, n := range nums {
		freq[n]++
	}
	pairs := []Pair{}
	for num, f := range freq {
		pairs = append(pairs, Pair{num, f})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].freq > pairs[j].freq
	})
	res := []int{}
	for i := 0; i < k; i++ {
		res = append(res, pairs[i].num)
	}
	return res
}

func main() {
	arr := []int{1, 1, 1, 2, 2, 3}
	k := 2
	fmt.Println(topKFrequent(arr, k))
}
