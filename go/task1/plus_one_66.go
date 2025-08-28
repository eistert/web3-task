package main

import "fmt"

/*
66. 加一
https://leetcode.cn/problems/plus-one/description/
*/
func plusOne(digits []int) []int {
	n := len(digits)

	for i := n - 1; i >= 0; i-- {
		if digits[i] != 9 {
			digits[i]++

			for j := i + 1; j < n; j++ {
				digits[j] = 0
			}

			return digits
		}
	}

	// digits 中所有的元素均为 9
	digits = make([]int, n+1)
	digits[0] = 1

	return digits
}

func main() {

	digits := []int{1, 2, 3}

	res := plusOne(digits)

	fmt.Println("res:", res)
}
