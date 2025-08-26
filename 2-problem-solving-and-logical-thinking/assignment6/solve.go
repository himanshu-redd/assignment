package main

import "fmt"

func rotate(matrix [][]int) {
	n := len(matrix)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	for i := 0; i < n; i++ {
		for l, r := 0, n-1; l < r; l, r = l+1, r-1 {
			matrix[i][l], matrix[i][r] = matrix[i][r], matrix[i][l]
		}
	}
}

func main() {
	m := [][]int{{1, 2}, {3, 4}}
	rotate(m)
	fmt.Println(m) // [[3 1] [4 2]]
}