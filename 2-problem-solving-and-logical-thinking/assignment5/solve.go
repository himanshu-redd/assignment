package main

import (
	"fmt"
)

func permute(str string) []string {
	var res []string
	chars := []rune(str)
	var backtrack func(start int)
	backtrack = func(start int) {
		if start == len(chars)-1 {
			res = append(res, string(chars))
			return
		}
		for i := start; i < len(chars); i++ {
			chars[start], chars[i] = chars[i], chars[start]
			backtrack(start + 1)
			chars[start], chars[i] = chars[i], chars[start]
		}
	}
	backtrack(0)
	return res
}

func main() {
	fmt.Println(permute("abc"))
}
