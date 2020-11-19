package v1

import (
	"fmt"
	"time"
)

type Cmd func(list []int) (ret []int)

//求偶数
func Evens(list []int) (ret []int) {
	ret = make([]int, 0)
	for _, num := range list {
		if num%2 == 0 {
			time.Sleep(time.Second * 1)
			ret = append(ret, num)
		}
	}
	return
}

func M10(list []int) (ret []int) {
	ret = make([]int, 0)
	for _, num := range list {
		time.Sleep(time.Millisecond * 300)
		ret = append(ret, num*10)
	}
	return
}

//管道函数
func Pipe(args []int, c1 Cmd, c2 Cmd) []int {
	ret := c1(args)
	return c2(ret)
}

func Test(nums []int) {
	ret := Pipe(nums, Evens, M10)
	for r := range ret {
		fmt.Printf("%d ", r)
	}
}
