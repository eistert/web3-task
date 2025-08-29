package main

import (
	"fmt"
	"sync"
)

/*
题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，
最后输出计数器的值。
考察点 ： sync.Mutex 的使用、并发数据安全。
*/

func main1() {
	var counter int       // 共享计数器
	var mu sync.Mutex     // 互斥锁
	var wg sync.WaitGroup // 等待组

	goroutines := 10
	increments := 1000

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				mu.Lock()   // 上锁
				counter++   // 修改共享变量
				mu.Unlock() // 解锁
			}
		}()
	}

	wg.Wait() // 等待所有 goroutine 执行完毕
	fmt.Println("最终计数器的值:", counter)

}
