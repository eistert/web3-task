package main

import "fmt"

/*
20. 有效的括号
考察：字符串处理、栈的使用
题目：给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
链接：https://leetcode-cn.com/problems/valid-parentheses/
*/

func isValid(s string) bool {
	n := len(s)
	if n%2 == 1 {
		return false
	}

	pairsMap := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	stack := []byte{}

	for i := 0; i < n; i++ {
		if pairsMap[s[i]] > 0 {
			if len(stack) == 0 || stack[len(stack)-1] != pairsMap[s[i]] {
				return false
			}

			stack = stack[:len(stack)-1]

		} else {
			stack = append(stack, s[i])
		}

	}

	return len(stack) == 0

}

func main() {

	s := "()"
	res := isValid(s)
	fmt.Println("res:", res)

}
