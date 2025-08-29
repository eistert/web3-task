package main

import (
	"fmt"
	"sync"
	"time"
)

/*
题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/

// 任务类型：返回 error 便于统计失败
type Task func() error

// 任务结果
type Result struct {
	Index    int
	Duration time.Duration
	Err      error
}

// RunTasks 并发执行所有任务，并返回每个任务的执行耗时与错误
func RunTasks(tasks []Task) []Result {
	n := len(tasks)
	results := make([]Result, n)

	var wg sync.WaitGroup
	wg.Add(n)

	for i, t := range tasks {
		i, t := i, t // 捕获循环变量
		go func() {
			defer wg.Done()

			start := time.Now()
			// 保护：如果任务内部 panic，也记录为错误
			defer func() {
				if r := recover(); r != nil {
					results[i] = Result{Index: i, Duration: time.Since(start), Err: fmt.Errorf("panic: %v", r)}
				}
			}()

			err := t()
			results[i] = Result{Index: i, Duration: time.Since(start), Err: err}
		}()
	}

	wg.Wait()
	return results
}

func main() {
	// 示例任务：用 Sleep 模拟不同耗时/错误
	tasks := []Task{
		func() error {
			time.Sleep(300 * time.Millisecond)
			return nil
		},
		func() error {
			time.Sleep(120 * time.Millisecond)
			return fmt.Errorf("something went wrong")
		},
		func() error {
			time.Sleep(50 * time.Millisecond)
			return nil
		},
		func() error { // 故意 panic 的任务
			time.Sleep(80 * time.Millisecond)
			panic("boom")
		},
	}

	results := RunTasks(tasks)

	// 输出结果
	for _, r := range results {
		if r.Err != nil {
			fmt.Printf("Task #%d: %v (duration=%v)\n", r.Index, r.Err, r.Duration)
		} else {
			fmt.Printf("Task #%d: OK (duration=%v)\n", r.Index, r.Duration)
		}
	}
}
