package main 

import "fmt"

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for i, num := range nums {
		if j, ok := m[target-num]; ok {
			return []int{j, i}
		}
		m[num] = i
	}
	return nil
}

func main() {
	arr := []int{2, 7, 11, 15}
	target := 9
	fmt.Println(twoSum(arr, target)) // [0, 1]
}
