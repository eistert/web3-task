package main

import "fmt"

/*
✅指针
题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
考察点 ：指针的使用、值传递与引用传递的区别。
*/

// 传引用：接收一个 *int
func addTenByPointer(num *int) {

	// *num = *num + 10
	// *num 解引用，取出指针指向的值，并加上 10。
	*num = *num + 10
	fmt.Println("在函数 addTenByPointer 内部:", *num)
}

// 传值：接收一个 int
func addTenByValue(num int) {
	num = num + 10
	fmt.Println("在函数 addTenByValue 内部:", num)
}

func main() {
	value := 20

	fmt.Println("初始值:", value)

	// --- 值传递 ---
	addTenByValue(value)
	fmt.Println("调用 addTenByValue 后:", value) // 不变，还是 20

	// --- 引用传递（指针） ---
	addTenByPointer(&value)
	fmt.Println("调用 addTenByPointer 后:", value) // 改变，变成 30
}
