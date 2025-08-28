package main

import "fmt"

/*
9. 回文数
https://leetcode.cn/problems/palindrome-number/description/
*/
func isPalindrome(x int) bool {
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}

	revertedNumber := 0

	for x > revertedNumber {
		revertedNumber = revertedNumber*10 + x%10
		x /= 10
	}

	b := x == revertedNumber
	b1 := x == revertedNumber/10

	return b || b1
}

func main() {
	x := 121

	res := isPalindrome(x)
	fmt.Println("res:", res)
}
