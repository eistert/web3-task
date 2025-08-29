package main

import "fmt"

func twoSum(nums []int, target int) []int {

	numsMap := make(map[int]int)

	for i, v := range nums {
		index, containKey := numsMap[target-v]

		if containKey {
			return []int{index, i}
		} else {
			numsMap[v] = i
		}
	}

	return []int{-1, -1}
}

func main5() {

	// arr := []int{2, 7, 11, 15}
	arr2 := []int{3, 2, 4}

	target := 6

	res := twoSum(arr2, target)

	fmt.Println("res:", res)

}
