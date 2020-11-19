package v2

import (
	"fmt"
	"time"
)

type Cmd func(list []int) chan int
type PipeCmd func(in chan int) chan int //支持管道的函数

//求偶数
func Evens(list []int) chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for _, num := range list {
			if num%2 == 0 { //业务流程
				time.Sleep(time.Second * 1)
				c <- num
			}
		}
	}()
	return c
}

//乘以2
func M2(in chan int) chan int { //这个函数是支持管道的
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			time.Sleep(time.Millisecond * 300)
			out <- num * 2
		}
	}()
	return out
}

//乘以10
func M10(in chan int) chan int { //这个函数是支持管道的
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			time.Sleep(time.Millisecond * 300)
			out <- num * 10
		}
	}()
	return out
}

//管道函数
func Pipe(args []int, c1 Cmd, cs ...PipeCmd) chan int {
	ret := c1(args)
	if len(cs) == 0 {
		return ret
	}
	retlist := make([]chan int, 0)
	for index, c := range cs {
		if index == 0 {
			retlist = append(retlist, c(ret))
		} else {
			getChan := retlist[len(retlist)-1]
			retlist = append(retlist, c(getChan))
		}
	}
	return retlist[len(retlist)-1]
}

func Test(nums []int) {
	ret := Pipe(nums, Evens, M10, M10, M10, M10)
	for r := range ret {
		fmt.Printf("%d ", r)
	}
}
