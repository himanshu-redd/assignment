package main

import (
	"fmt"
)

func minWindow(s string, t string) string {
	if len(t) == 0 || len(s) == 0 {
		return ""
	}

	need := make(map[byte]int)
	for i := range t {
		need[t[i]]++
	}

	have := make(map[byte]int)
	required := len(need)
	formed := 0
	l, r := 0, 0
	ans := []int{-1, 0, 0} // length, left, right

	for r < len(s) {
		c := s[r]
		have[c]++
		if need[c] > 0 && have[c] == need[c] {
			formed++
		}

		for l <= r && formed == required {
			if ans[0] == -1 || r-l+1 < ans[0] {
				ans[0] = r - l + 1
				ans[1] = l
				ans[2] = r
			}
			ch := s[l]
			have[ch]--
			if need[ch] > 0 && have[ch] < need[ch] {
				formed--
			}
			l++
		}
		r++
	}

	if ans[0] == -1 {
		return ""
	}
	return s[ans[1] : ans[2]+1]
}

func main() {
	s := "ADOBECODEBANC"
	t := "ABC"
	fmt.Println(minWindow(s, t)) // Output: "BANC"
}
