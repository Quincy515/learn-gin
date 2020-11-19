Go 高并发模型之管道模式

[toc]

### 01. 认识管道模式、最初级的代码

例子：Linux命令很经典的命令  `cat log.txt | grep 'abc'`

基本概念是：前面一个进程的输出（`stdout`）直接作为下一个进程的输入（`stdin`）。

可以使用程序实现**类似的模式**。

#### 从最简单的开始

假设有一个切片 `list:=[]int{2,3,6,12,22,16,4,9,23,64,62}`

有两步需求:

1、从里面找到偶数    

2、把偶数乘以10

```go
package main

import (
	"fmt"
)

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

func main() {
	nums := []int{2, 3, 6, 12, 22, 16, 4, 9, 23, 64, 62}
	fmt.Println(Multiply(Evens(nums))) // 函数的嵌套调用
}
```

如果换成 `Linux` 命令应该怎么敲 `Evens nums | Multiply`

最初级的封装

```go
func p(args []int,c1 Cmd,c2 Cmd) []int {
	 ret:= c1(args)
	 return c2(ret)
}
```

完整代码

```go
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
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/992e86154ef02f518a640b5f3ee72d53d3111a06#diff-e695b9090228290d541dcaaff7ccf49dc68debd598c1f2520f8d03afe27f98f7R1)

### 02. 使用 channel 改进管道模式

约定：

> 凡是支持管道模式的函数，其参数必须是 channel。返回 channel。

```go
package main

import "fmt"

type Cmd func(list []int) chan int
type PipeCmd func(in chan int) chan int //支持管道的函数

// Evens 求偶数
func Evens(list []int) chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for _, num := range list {
			if num%2 == 0 { //业务流程
				c <- num
			}
		}
	}()

	return c

}

// M10 乘以10
func M10(in chan int) chan int { //这个函数是支持管道的
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			out <- num * 10
		}
	}()
	return out
}

//管道函数
func Pipe(args []int, c1 Cmd, c2 PipeCmd) chan int {
	ret := c1(args)
	return c2(ret)
}

func main() {
	nums := []int{2, 3, 6, 12, 22, 16, 4, 9, 23, 64, 62}

	ret := Pipe(nums, Evens, M10)
	for r := range ret {
		fmt.Printf("%d ", r)
	}
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/2949e84ca7cadf3cbc73fac5a0f0954c9c311f95#diff-e695b9090228290d541dcaaff7ccf49dc68debd598c1f2520f8d03afe27f98f7L2)

### 03. 管道模式性能对比，可变参数

代码变动 [git commit]()

