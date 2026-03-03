package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// var num int8 = -128                           // 十进制数
	// binaryStr := strconv.FormatInt(int64(num), 2) // 将int64(num)转换为二进制字符串
	// fmt.Println("Binary:", binaryStr)
	//slice := []int{1, 2, 3, 4}
	slice := make([]int, 5, 10)
	fmt.Println("输出 slice:", unsafe.Sizeof(slice), cap(slice), len(slice), slice)
	slice = append(slice, 5, 6, 7, 8, 9, 10, 11)
	fmt.Println("输出 slice:", unsafe.Sizeof(slice), cap(slice), len(slice), slice)
	slice = append(slice, 5, 6, 7, 8, 9, 10, 11)
	fmt.Println("输出 slice:", unsafe.Sizeof(slice), cap(slice), len(slice), slice)

	// 基本for循环
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// 类似while循环
	i := 0
	for i < 5 {
		fmt.Println(i)
		i++
	}
	i1 := 0
	// 死循环
	for {
		// 需要在循环体内使用break跳出
		if i1 > 5 {
			fmt.Println(i1)
			break
		}
		i1++
	}

	// range遍历
	arr := []int{1, 2, 3, 4, 5}
	for index, value := range arr {
		fmt.Printf("index: %d, value: %d\n", index, value)
	}

	// 只遍历值
	for _, value := range arr {
		fmt.Println(value)
	}

	// 遍历字符串（按rune）
	str := "Hello 世界"
	for _, char := range str {
		fmt.Printf("%c ", char)
	}

	// 遍历映射
	m := map[string]int{"a": 1, "b": 2}
	for key, value := range m {
		fmt.Printf("%s: %d\n", key, value)
	}

	fmt.Printf("结束了")
}
