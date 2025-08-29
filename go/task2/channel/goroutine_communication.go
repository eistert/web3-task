package main

import (
	"fmt"
	"sync"
)

/*
题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，
另一个协程从通道中接收这些整数并打印出来。
考察点 ：通道的基本使用、协程间通信。
*/

func main1() {

	// 无缓冲通道
	ch := make(chan int)

	var wg sync.WaitGroup
	wg.Add(1)

	/*
		这是 匿名函数 + 立即调用 的写法。

		func() { ... } → 定义了一个匿名函数

		go func() { ... }() → 启动一个 goroutine 来执行它

		最后的 () 就是“立即调用”
	*/
	// 生产者：1..10 写入通道，结束后关闭通道
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch) // very important：通知接收方没有更多数据
	}()

	// 消费者：从通道读取并打印，直到通道被关闭
	go func() {
		defer wg.Done()
		for v := range ch { // ch 关闭后自动退出循环
			fmt.Println(v)
		}
	}()

	wg.Wait()
	fmt.Println("完成")
}
