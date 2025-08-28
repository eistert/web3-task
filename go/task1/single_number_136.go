package main

import "fmt"

/*
 136. 只出现一次的数字

https://leetcode.cn/problems/single-number/description/
*/
func singleNumber(nums []int) int {
	countmap := make(map[int]int)

	// 使用range遍历
	for _, v := range nums {
		countmap[v] = countmap[v] + 1 // 或者直接 countmap[v]++
	}

	// 遍历map集合main
	for k, v := range countmap {
		if v == 1 {
			return k
		}
	}

	return 0
}

// 利用 a^a=0、a^0=a
func singleNumberXOR(nums []int) int {
	res := 0
	for _, x := range nums {
		res ^= x
	}
	return res
}

func main() {

	var b []int = []int{4, 1, 2, 1, 2}
	res := singleNumber(b)
	fmt.Printf("res:%d\n", res)
}
