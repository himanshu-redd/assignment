package main

import "fmt"

func main() {
	arr1 := []int{1, 2, 3, 4, 5, 6}
	arr2 := []int{3, 5, 7, 9}

	i, j := 0, 0
	fmt.Println("Intersection:")
	for i < len(arr1) && j < len(arr2) {
		if arr1[i] == arr2[j] {
			fmt.Println(arr1[i])
			i++
			j++
		} else if arr1[i] < arr2[j] {
			i++
		} else {
			j++
		}
	}
}