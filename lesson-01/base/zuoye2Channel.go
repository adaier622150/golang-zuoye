package main

import (
	"fmt"
	"time"
)

// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，
//另一个协程从通道中接收这些整数并打印出来。
//考察点 ：通道的基本使用、协程间通信。
//题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
//考察点 ：通道的缓冲机制。

func chan1(ch chan int) {
	defer fmt.Println("程序结束 chan1")
	go func() {
		defer close(ch)
		defer fmt.Println("close ch")
		for i := 0; i < 10; i++ {
			ch <- i
			fmt.Println("写入通道数据 ", i)
		}
	}()
}
func chan2(ch chan int) {
	defer fmt.Println("程序结束 chan2")
	go func() {
		for v := range ch {
			fmt.Println("读取通道数据", v)
			time.Sleep(100 * time.Millisecond)
		}
	}()
}
func producer(ch chan int) {
	defer fmt.Println("程序结束 producer")
	go func() {
		defer close(ch)
		defer fmt.Println("close ch")
		for i := 0; i < 100; i++ {
			ch <- i
			fmt.Println("写入通道数据 ", i)
		}
	}()
}
func consumer(ch chan int) {
	defer fmt.Println("程序结束 consumer")
	go func() {
		for v := range ch {
			fmt.Println("读取通道数据", v)
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func main() {
	defer fmt.Println("程序结束")
	ch := make(chan int, 10)
	go producer(ch)
	go consumer(ch)
	fmt.Println("等待goroutine完成,休息5秒")
	time.Sleep(10 * time.Second)
}
