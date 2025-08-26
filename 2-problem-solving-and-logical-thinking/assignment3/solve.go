package main 

import (
	"fmt"
	"sort"
	"strings"
)

func isAnagram(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	s1Slice := strings.Split(s1, "")
	s2Slice := strings.Split(s2, "")

	sort.Strings(s1Slice)
	sort.Strings(s2Slice)

	for i := range s1Slice {
		if s1Slice[i] != s2Slice[i] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(isAnagram("listen", "silent")) // true
	fmt.Println(isAnagram("hello", "world"))   // false
}
