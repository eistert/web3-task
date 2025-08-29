package main

import "fmt"

/*
题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/

/*
额外说明：在实际开发中，因为 切片本身是引用类型，一般直接传 []int 就能修改里面的元素，
不需要再传 *[]int。不过这道题是专门考“切片指针”，所以用了 *[]int 的写法。
*/
func sliceInt(slice []int) {
	for i, v := range slice {
		slice[i] = v * 2
	}
}

// 函数：接收一个整数切片的指针，将每个元素 *2
func doubleSlice(nums *[]int) {
	for i := 0; i < len(*nums); i++ {
		(*nums)[i] = (*nums)[i] * 2
	}
}

func main() {

	s1 := []int{5, 4, 3, 2, 1}
	fmt.Println("修改前 s1:", s1)

	sliceInt(s1)

	fmt.Println("修改后 s1:", s1)

	// // 定义一个切片
	// values := []int{1, 2, 3, 4, 5}
	// fmt.Println("修改前:", values)

	// // 传入切片指针
	// doubleSlice(&values)

	// fmt.Println("修改后:", values)
}
