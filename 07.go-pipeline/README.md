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

目录 `v1` 是没有使用 `channel` ，目录 `v2` 是使用 `channel`。

模拟管道的函数 `Pipe` 修改为可变参数

```go
//管道函数
func Pipe(args []int, c1 Cmd, cs ...PipeCmd) chan int {
	ret := c1(args)
	if len(cs) == 0 {
		return ret
	}
	retlist := make([]chan int, 0)
	for index, c := range cs {
		if index == 0 { // 第一次执行
			retlist = append(retlist, c(ret)) // 第一个执行结果放入切片
		} else { // 第二次执行，获取结果切片中最后一个作为结果
			getChan := retlist[len(retlist)-1]
			retlist = append(retlist, c(getChan))
		}
	}
	return retlist[len(retlist)-1] // 返回结果切片中最后一个结果
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/e7dd79a999c07337d4eafef86233321a22c8e6c7#diff-e695b9090228290d541dcaaff7ccf49dc68debd598c1f2520f8d03afe27f98f7L1)

### 04. 管道模式之多路复用、提高性能

上面实现的程序，执行顺序是：**找偶数** => **偶数乘以10** => **偶数乘以2**

使用了 `channel` 之后的顺序：**找偶数** => **偶数乘以10** => **偶数乘以2**

-----------只要完成一个操作-----------↓-------- ↗----- ↓---------↗

-------------就放入channel---------channel -------channel

多路复用：               **找偶数**

----------------------------- ↓

-------------------------channel

---------------------↙ ----------- ↘

----------------**偶数乘以10** ----**偶数乘以10** 

多个同样的函数，同时对 `channel` 进行读值。

> 多个函数同时从同一个channel里读取数据。直至channel被关闭
>
> 可以更好的利用多核。

注意要把之前求偶数的函数 `Even()` 中，模拟等待操作时间注释掉

```go
//求偶数
func Evens(list []int) chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for _, num := range list {
			if num%2 == 0 { //业务流程
				//time.Sleep(time.Second * 1)
				c <- num
			}
		}
	}()
	return c
}
```

原来的管道函数

```go
// Pipe 管道函数
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
```

多路复用实现

```go
// Pipe2 多了复用
func Pipe2(args []int, c1 Cmd, cs ...PipeCmd) chan int {
	ret := c1(args)
	out := make(chan int)
	wg := sync.WaitGroup{}
	for _, c := range cs {
		getChan := c(ret)
		wg.Add(1)
		go func(input chan int) {
			defer wg.Done()
			for v := range input {
				out <- v
			}
		}(getChan)
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()
	return out
}
```

 代码变动 [git commit]()