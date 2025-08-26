package main

import "fmt"

func partition(nums []int, left, right int) int {
	pivot := nums[right]
	i := left
	for j := left; j < right; j++ {
		if nums[j] > pivot { // note: > for kth largest
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
	}
	nums[i], nums[right] = nums[right], nums[i]
	return i
}

func quickSelect(nums []int, left, right, k int) int {
	if left <= right {
		p := partition(nums, left, right)
		if p == k {
			return nums[p]
		} else if p < k {
			return quickSelect(nums, p+1, right, k)
		} else {
			return quickSelect(nums, left, p-1, k)
		}
	}
	return -1
}

func findKthLargestQuickSelect(nums []int, k int) int {
	return quickSelect(nums, 0, len(nums)-1, k-1)
}

func main() {
	arr := []int{3, 2, 1, 5, 6, 4}
	k := 5 
	fmt.Println(findKthLargestQuickSelect(arr, k)) // Output: 5
}
