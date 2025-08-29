package main

import (
	"fmt"
	"math"
)

/*
题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
考察点 ：接口的定义与实现、面向对象编程风格。
*/

// 定义 Shape 接口
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 定义矩形结构体
type Rectangle struct {
	Width, Height float64
}

// 实现 Shape 接口的方法
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 定义圆形结构体
type Circle struct {
	Radius float64
}

// 实现 Shape 接口的方法
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func main1() {
	// 创建实例
	rect := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 4}

	// 定义一个 Shape 类型的切片，存放不同形状
	shapes := []Shape{rect, circle}

	// 遍历并调用接口方法
	for _, s := range shapes {
		fmt.Printf("类型: %T\n", s)
		fmt.Printf("面积: %.2f\n", s.Area())
		fmt.Printf("周长: %.2f\n\n", s.Perimeter())
	}
}
