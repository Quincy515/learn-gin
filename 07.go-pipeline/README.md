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

代码变动 [git commit]()