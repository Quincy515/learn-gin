package main

import "fmt"

type Cmd func(list []int) (ret []int)

// Evens 寻找偶数
func Evens(list []int) (ret []int) {
	ret = make([]int, 0)
	for _, num := range list {
		if num%2 == 0 {
			ret = append(ret, num)
		}
	}
	return
}

// Multiply 乘以 10
func Multiply(list []int) (ret []int) {
	ret = make([]int, 0)
	for _, num := range list {
		ret = append(ret, num*10)
	}
	return
}

// p 模拟管道函数
func p(args []int, c1 Cmd, c2 Cmd) []int {
	ret := c1(args)
	return c2(ret)
}

func main() {
	nums := []int{2, 3, 6, 12, 22, 16, 4, 9, 23, 64, 62}
	//fmt.Println(Multiply(Evens(nums)))
	fmt.Println(p(nums, Evens, Multiply))
}
