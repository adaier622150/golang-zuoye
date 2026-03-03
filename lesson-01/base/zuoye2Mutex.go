package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
//考察点 ： sync.Mutex 的使用、并发数据安全。
//题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
//考察点 ：原子操作、并发数据安全。

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (sc *SafeCounter) Increment() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.count++
}

func (sc *SafeCounter) getCount() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.count
}
func main() {
	var wg sync.WaitGroup
	//sc := &SafeCounter{}
	mu := sync.Mutex{}
	sum := 0
	fmt.Println("加锁 开始")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 100; j++ {
				//sc.Increment()
				mu.Lock()
				sum++
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("加锁 计数器值:", sum)
	//fmt.Println("加锁 计数器值:", sc.getCount())
	fmt.Println("无锁 开始")
	var counter int64
	//var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 100; j++ {
				atomic.AddInt64(&counter, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("无锁 计数器值:", counter)
}
