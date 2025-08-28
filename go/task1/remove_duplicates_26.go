package main

/*
总结
fast：扫描数组找新元素
slow：记录下一个存放位置
当遇到新元素时：把它放到 nums[slow]，并 slow++
最终 slow 就是去重后数组的长度
*/
func removeDuplicates(nums []int) int {
	n := len(nums)

	if n == 0 {
		return 0
	}

	slow := 1

	for fast := 1; fast < n; fast++ {
		if nums[fast] != nums[fast-1] {
			nums[slow] = nums[fast]
			slow++
		}
	}

	return slow
}
