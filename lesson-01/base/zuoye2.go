package main

import (
	"fmt"
	"time"
)

// 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
// 考察点 ：指针的使用、值传递与引用传递的区别。
// 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
// 考察点 ：指针运算、切片操作。
func one1(a *int) {
	*a = *a + 10
}

func one2(nums *[]int) {
	fmt.Println("nums ", nums)
	for index, value := range *nums {
		fmt.Println("i ", index, value)
		(*nums)[index] = value * 2
	}
}
func modifyArray(arrPtr *[5]int) {
	for i := range arrPtr {
		arrPtr[i] *= 2 // 将每个元素乘以2
	}
}
func multiplyByTwo(slice *[]int) {
	for i := range *slice {
		(*slice)[i] *= 2
	}
}

// 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
// 考察点 ： go 关键字的使用、协程的并发执行。
// 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 考察点 ：协程原理、并发任务调度。
// zuoye2Schedule.go
func run1() {
	go goroutine1()
	go goroutine2()
	fmt.Println("等待goroutine完成,休息5秒")
	time.Sleep(5 * time.Second)
}

func goroutine1() {
	defer fmt.Println("奇数 程序结束")
	for i := 1; i <= 10; i += 2 {
		fmt.Println("奇数：", i)
		time.Sleep(100 * time.Millisecond)
	}
}

func goroutine2() {
	defer fmt.Println("偶数 程序结束")
	for i := 2; i <= 10; i += 2 {
		fmt.Println("偶数：", i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	//nums := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	//nums := []int{2, 7, 11, 15}
	//nums1 := [5]int{2, 7, 11, 15}
	//target := 22
	//fmt.Println("合并区间", nums)
	//x := TwoSum(nums, target)
	//fmt.Println("合并区间", x)

}
