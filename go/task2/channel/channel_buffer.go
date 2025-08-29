package main

import (
	"fmt"
	"sync"
)

/*
题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
考察点 ：通道的缓冲机制。
*/

func main() {

	ch := make(chan int, 10)

	var wg sync.WaitGroup
	wg.Add(1)

	// 生产者
	go func() {
		for i := 0; i < 100; i++ {
			ch <- i // 如果缓冲满了，会阻塞直到有空间
		}
		close(ch) // 发送完毕后关闭通道
	}()

	// 消费者
	go func() {
		defer wg.Done()
		for v := range ch { // 从通道接收，直到通道被关闭
			fmt.Println("消费者:", v)
		}
	}()

	wg.Wait()
	fmt.Println("所有数据处理完毕")

}
