package main
import (
	"fmt"
)

func maxSlidingWindow(nums []int, k int) []int {
	var result []int
	deque := []int{} 

	for i := 0; i < len(nums); i++ {
		for len(deque) > 0 && nums[deque[len(deque)-1]] < nums[i] {
			deque = deque[:len(deque)-1]
		}

		deque = append(deque, i)

		if deque[0] <= i-k {
			deque = deque[1:]
		}

		if i >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}

	return result
}

func main() {
	arr := []int{1, 3, -1, -3, 5, 3, 6, 7}
	k := 3
	fmt.Println(maxSlidingWindow(arr, k)) // [3 3 5 5 6 7]
}
