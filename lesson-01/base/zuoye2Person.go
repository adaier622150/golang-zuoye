package main

import (
	"fmt"
)

// 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，
// 组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
// 考察点 ：组合的使用、方法接收者。

type Person struct {
	Name string
	Age  int
}

func (p Person) GetInfo() string {
	return fmt.Sprintf("员工西悉尼 %s,年龄 %d", p.Name, p.Age)
}

type Employee struct {
	Person
	EmployeeID string
}

func (e Employee) PrintInfo() {
	fmt.Printf("员工信息：%s,   %s\n", e.GetInfo(), e.EmployeeID)
}

// PrintShapeInfo 多态函数

func main() {

	employee := Employee{
		Person: Person{
			Name: "张三",
			Age:  30,
		},
		EmployeeID: "E001",
	}
	employee.PrintInfo()
}
