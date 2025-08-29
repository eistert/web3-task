package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter int64 // 共享计数器
	var wg sync.WaitGroup
	goroutines := 10
	increments := 1000

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				atomic.AddInt64(&counter, 1) // 原子递增
			}
		}()
	}

	wg.Wait()
	fmt.Println("最终计数器的值:", atomic.LoadInt64(&counter)) // 原子读取

}
