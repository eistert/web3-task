package main

import "fmt"

/*
题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。
为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
考察点 ：组合的使用、方法接收者。
*/

// 定义 Person 结构体
type Person struct {
	Name string
	Age  int
}

// 定义 Employee 结构体，组合 Person
type Employee struct {
	Person
	EmployeeID string
}

func (e Employee) PrintInfo() {
	fmt.Printf("员工姓名: %s\n", e.Name)
	fmt.Printf("员工年龄: %d\n", e.Age)
	fmt.Printf("员工ID: %s\n", e.EmployeeID)
}

func main() {

	// 创建一个 Employee 实例
	emp := Employee{
		Person: Person{
			Name: "张三",
			Age:  30,
		},
		EmployeeID: "E1001",
	}

	emp.PrintInfo()

}
