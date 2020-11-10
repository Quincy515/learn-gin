gin 脚手架研发

### 01. 从零开始

新建目录 `src/cmd/main.go`

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()
	r.Handle("GET", "/", func(c *gin.Context){
		c.JSON(200, gin.H{"result":"success"})
	})

	r.Run(":8080")
}
```

请求 http://localhost:8080/ 可以看到 `{ "result": "success" }`

### 02. 隐藏路由和业务方法

```go
package main

import (
	. "gin-up/src/classes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	NewIndexClass(r).Build() // 路由和业务方法隐藏

	r.Run(":8080")
}
```

在 `src` 目录下新建 `classes` 目录和 `index.go` 文件

```go
package classes

import "github.com/gin-gonic/gin"

// IndexClass 嵌套 *gin.Engine
type IndexClass struct {
	*gin.Engine // gin.New() 创建的
	// 嵌套，好比继承，但不是继承
}

// NewIndexClass 所谓的构造函数
func NewIndexClass(e *gin.Engine) *IndexClass {
	return &IndexClass{Engine: e} // 需要赋值，因为是指针
}

// GetIndex 业务方法，函数名根据业务而起
func (i *IndexClass) GetIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": "index ok",
		})
	}
}

// Build 把业务的路由隐藏在 Build 函数
func (i *IndexClass) Build() {
	i.Handle("GET", "/", i.GetIndex())
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/d103180c4c866f505adcbc9b9b367b26dda12397#diff-5f86647f5f70db405f26a54014d77b1d15d913f96b612dc6763b8870041577d8R1)

### 03. 自定义快捷模板

在 `src/classes` 目录下新建 `User.go` 文件，注意**文件名统一首字母大写** 以便于使用快捷模板

```go
func (this *$FileName$Class) Build() {
	this.Handle("GET", "/$path$", this.FuncName())
}
```

<img src="../imgs/18_live_templates.png" style="zoom:95%;" />

```go
func (this *$FileName$Class) FuncName() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": "success",
		})
	}
}
```

<img src="../imgs/19_live_templates.png" style="zoom:95%;" />

这样在 `User.go` 中，手写 

```go
package classes

import "github.com/gin-gonic/gin"

// UserClass *gin.Engine 的嵌套
type UserClass struct {
	*gin.Engine
}
```

鼠标放到 `type UserClass struct {` 右键选择 `generate` 再选择 `constructor` 自动生成 **构造函数**。

```go
// NewUserClass UserClass generate constructor
func NewUserClass(engine *gin.Engine) *UserClass {
	return &UserClass{Engine: engine}
}
```

然后使用快捷键盘 `method` 自动生成 **控制器函数**

```go
// UserList 控制器方法
func (this *UserClass) UserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": "success",
		})
	}
}
```

继续使用快捷键 `build` 自动生成 **控制器 build 方法**

```go
func (this *UserClass) Build() {
	this.Handle("GET", "/user", this.UserList())
}
```

最后在 `main.go` 中调用 `NewUserClass(r).Build()`

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/64c4800eb89dc38914dca1c1f3aba11a1b880187#diff-8d9e1f78703b2eb32787b5d6fcdc6da3201ad241fb4c572b6bbe8eb8284031e3R1)

### 04. 封装 gin 成自己的框架 goft

现在 `main.go` 的代码形式是

```go
func main() {
	r := gin.New()
	NewIndexClass(r).Build() // 路由和业务方法隐藏
	NewUserClass(r).Build() // 控制器
	r.Run(":8080")
}
```

还有一些问题：

1. 如果控制器多了，代码还是冗余
2. 各个控制器代码之前，没有约束(没有接口规范)

新建目录 `src\goft` 和文件 `Goft.go`

```go
package goft

import "github.com/gin-gonic/gin"

// Goft 嵌套 *gin.Engine
type Goft struct {
	*gin.Engine
}

// Ignite Goft 的构造函数，发射、燃烧，富含激情的意思
func Ignite() *Goft {
	return &Goft{Engine: gin.New()}
}

// Launch 最终启动函数，相当于 r.Run()
func (this *Goft) Launch() {
	this.Run(":8080")
}

// Mount 挂载控制器，定义接口，接口里的方法作为参数，控制器实现接口就可以传进来
func (this *Goft) Mount() {

}
```

在 `src/goft` 目录下新建 `IClass.go` 文件，定义接口

```go
package goft

type IClass struct {
	Build(goft *Goft)
}
```

修改业务控制器 `src/classes/User.go` 代码，删除 `*gin.Engine`，

使得业务控制器和服务器没有强关联，否则耦合太高。

```go
// UserClass 
type UserClass struct {}

// NewUserClass UserClass generate constructor
func NewUserClass() *UserClass {
	return &UserClass{}
}
```

控制器生成方法 `Build()` 修改传入参数，这样就实现了 `IClass` 接口

```go
func (this *UserClass) Build(goft *goft.Goft) {
	goft.Handle("GET", "/user", this.UserList())
}
```

在 `Mount()` 函数里就可以传入参数

```go
func (this *Goft) Mount(classes ...IClass) *Goft {
	for _, class := range classes {
		class.Build(this) 
	}
    return this
}
```

返回自己 `*Goft` 是为了链式调用。同样的修改 `src/classes/Index.go`

```go
package classes

import (
	"gin-up/src/goft"
	"github.com/gin-gonic/gin"
)

type IndexClass struct{}

func NewIndexClass() *IndexClass {
	return &IndexClass{}
}

// GetIndex 业务方法，函数名根据业务而起
func (i *IndexClass) GetIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": "index ok",
		})
	}
}

// Build 把业务的路由隐藏在 Build 函数
func (i *IndexClass) Build(goft *goft.Goft) {
	goft.Handle("GET", "/", i.GetIndex())
}
```

在 `main.go` 中调用

```go
package main

import (
	. "gin-up/src/classes"
	"gin-up/src/goft"
)

func main() {
	goft.Ignite().Mount(NewIndexClass(), NewUserClass()).Launch()
}
```

最终的 `src/goft/Goft.go` 代码

```go
package goft

import "github.com/gin-gonic/gin"

// Goft 嵌套 *gin.Engine
type Goft struct {
	*gin.Engine
}

// Ignite Goft 的构造函数，发射、燃烧，富含激情的意思
func Ignite() *Goft {
	return &Goft{Engine: gin.New()}
}

// Launch 最终启动函数，相当于 r.Run()
func (this *Goft) Launch() {
	this.Run(":8080")
}

// Mount 挂载控制器，定义接口，控制器继承接口就可以传进来
func (this *Goft) Mount(classes ...IClass) *Goft {
	for _, class := range classes {
		class.Build(this)
	}
	return this
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/42141001ea2945a28c70beaaf3f9560761a59e94#diff-8d9e1f78703b2eb32787b5d6fcdc6da3201ad241fb4c572b6bbe8eb8284031e3L1)

### 05. 把路由挂载到 Group 中

`gin` 可以自定义 `group` 分组，方便进行 API 版本的管理

```go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	//goft.Ignite().Mount(NewIndexClass(), NewUserClass()).Launch()

	r := gin.New()
	v1 := r.Group("v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}
	r.Run()
}
```

在 `Mount()` 函数中修改代码实现 `Group` 功能

```go
func (this *Goft) Mount(group string, classes ...IClass) *Goft {
	this.g = this.Group(group)
	for _, class := range classes {
		class.Build(this)
	}
	return this
}
```

修改主类 `Goft`

```go
type Goft struct {
	*gin.Engine // 把 *gin.Engine 放入主类里
	g *gin.RouterGroup // 保存 group 对象
}
```

不改变控制器代码，可以重载 `gin.Handle()`  函数

```go
// Handle 重载 gin.Handle 函数
func (this *Goft) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) *Goft {
	this.g.Handle(httpMethod, relativePath, handlers...) // 最后一个参数，需要使用 ...来延展
	return this
}
```

这样在 `main.go` 中就可以加入 `group` 参数

```go
package main

import (
	. "gin-up/src/classes"
	"gin-up/src/goft"
)

func main() {
	goft.Ignite().
		Mount("v1", NewIndexClass(),
			NewUserClass()).
		Mount("v2", NewIndexClass()).
		Launch()
}
```

代码修改 [git commit]()