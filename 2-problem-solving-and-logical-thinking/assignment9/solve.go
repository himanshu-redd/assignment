package main

import "fmt"

func trap(height []int) int {
    l, r := 0, len(height)-1
    lmax, rmax, res := 0, 0, 0
    for l < r {
        if height[l] < height[r] {
            if height[l] >= lmax {
                lmax = height[l]
            } else {
                res += lmax - height[l]
            }
            l++
        } else {
            if height[r] >= rmax {
                rmax = height[r]
            } else {
                res += rmax - height[r]
            }
            r--
        }
    }
    return res
}

func main() {
    fmt.Println(trap([]int{0,1,0,2,1,0,1,3,2,1,2,1}))
}
