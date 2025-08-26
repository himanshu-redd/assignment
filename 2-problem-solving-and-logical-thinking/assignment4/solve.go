package main 

import "fmt"

func isValid(s string) bool {
	stack := []rune{}
	mapping := map[rune]rune{')': '(', ']': '[', '}': '{'}

	for _, ch := range s {
		if ch == '(' || ch == '[' || ch == '{' {
			stack = append(stack, ch)
		} else {
			if len(stack) == 0 || stack[len(stack)-1] != mapping[ch] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func main() {
	fmt.Println(isValid("{[()]}")) // true
	fmt.Println(isValid("{[(])}")) // false
	fmt.Println(isValid("()[]{}")) // true
}
