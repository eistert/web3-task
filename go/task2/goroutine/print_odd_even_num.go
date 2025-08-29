package main

import (
	"fmt"
	"sync"
)

/*
题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
*/

func main1() {

	var wg sync.WaitGroup
	wg.Add(2) // 等待两个协程

	printOdd(&wg)
	printEven(&wg)

	// 等待两个协程完成
	wg.Wait()
	fmt.Println("全部打印完成")

}

func printOdd(wg *sync.WaitGroup) {
	defer wg.Done() // 协程结束时计数减一
	for i := 1; i <= 10; i += 2 {
		fmt.Println("奇数:", i)
	}

}

func printEven(wg *sync.WaitGroup) {
	defer wg.Done() // 协程结束时计数减一
	for i := 2; i <= 10; i += 2 {
		fmt.Println("偶数:", i)
	}
}
