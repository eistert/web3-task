package main

import (
	"fmt"
	"sort"
)

func merge(intervals [][]int) [][]int {

	n := len(intervals)

	if n == 0 {
		return nil
	}

	// 1) 按起点升序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// 2) 线性扫描合并
	res := make([][]int, 0, n)

	for _, cur := range intervals {
		// 如果结果为空，或当前区间与最后区间不重叠，直接加入
		if len(res) == 0 || cur[0] > res[len(res)-1][1] {
			// 拷贝一份，避免后续修改影响原切片
			res = append(res, []int{cur[0], cur[1]})
		} else {
			// 否则合并到最后一个区间：更新终点为较大者
			last := res[len(res)-1]
			if cur[1] > last[1] {
				last[1] = cur[1]
			}
		}
	}

	return res
}

func main1() {
	input := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Println(merge(input)) // [[1 6] [8 10] [15 18]]
}
