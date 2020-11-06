gin 开发代码技巧

[toc]

### 01. 实体定义技巧

最基础代码

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		
	})
}
```

其中 `func(c *gin.Context) {}` 可以封装成 `controller` 控制器单独处理业务。

比如最基本的实体定义代码如下

```go
package main

import "github.com/gin-gonic/gin"

type UserModel struct {
	UserID   int
	UserName string
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func main() {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		user := NewUserModel()
		c.JSON(200, user)
	})
}
```

可以定义一个 `src` 目录，在 `src` 目录下再新建一个 `models` 目录，把所有的模型放入 `models` 里。

这样就可以把 `UserModel` 放入 `models` 文件夹里，这里建议再新建一个 `UserModel` 文件夹，

使用时可以 `package UserModel` 方便使用。

在 `UserModel` 文件夹里随便新建一个文件 `model.go`

```go
package UserModel

type UserModelImpl struct {
	UserID   int
	UserName string
}

func New() *UserModelImpl {
	return &UserModelImpl{}
}
```

这样再使用时 `UserModel.New()` 比较方便易懂

```go
package main

import (
	"ginskill/src/models/UserModel"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		user := UserModel.New()
		c.JSON(200, user)
	})
	r.Run(":8080")
}
```

此时文件目录如下

```bash
└─src
    └─models
        └─UserModel
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/4c7bc60999f5ca8330e8f3041f06c1e91c404151#diff-2357d785d351a1c8beb39645ad84efb5d68e67b935ba83881b80a1e329ba2c64R1)

### 02. 实体带参数的初始化技巧

上面目录的好处就是借助包名 `UserModel.New()` 完成初始化，

如果希望对里面的属性进行赋值，有以下几种方式

#### 方案一：重写

```go
package UserModel

type UserModelImpl struct {
	UserID   int
	UserName string
}

func New() *UserModelImpl {
	return &UserModelImpl{}
}

func NewWithID(id int) *UserModelImpl {
	return &UserModelImpl{UserID: id}
}

func NewWithName(name string) *UserModelImpl {
	return &UserModelImpl{UserName: name}
}
```

#### 方案二：可变参数

新建 `attrs.go` 文件

```go
package UserModel

type UserModelAttrFunc func(u *UserModelImpl)

type UserModelAttrFuncs []UserModelAttrFunc

func (this UserModelAttrFuncs) Apply(u *UserModelImpl) {
	for _, f := range this {
		f(u)
	}
}
```

然后修改 `model.go` 文件为可变参数

```go
package UserModel

type UserModelImpl struct {
	UserID   int
	UserName string
}

func New(attrs ...UserModelAttrFunc) *UserModelImpl {
	u := &UserModelImpl{}
	// 对 u 里每个属性进行初始化
	// 强制类型转化。
	UserModelAttrFuncs(attrs).Apply(u)
	return u
}
```

这样在 `attrs.go` 中就可以对每个属性执行不同的初始化

```go
package UserModel

type UserModelAttrFunc func(u *UserModelImpl)

type UserModelAttrFuncs []UserModelAttrFunc

func WithUserID(id int) UserModelAttrFunc {
	return func(u *UserModelImpl) {
		u.UserID = id
	}
}
func WithUserName(name string) UserModelAttrFunc {
	return func(u *UserModelImpl) {
		u.UserName = name
	}
}

func (this UserModelAttrFuncs) Apply(u *UserModelImpl) {
	for _, f := range this {
		f(u)
	}
}
```

在初始化时 `user := UserModel.New()` 默认为空的初始化，实体实例化

`user := UserModel.New(UserModel.WithUserID(101))`

`user := UserModel.New(UserModel.WithUserID(101), UserModel.WithUserName("custer"))`

### 03. 链式调用

比如 `user := UserModel.New(UserModel.WithUserID(101))`  初始化之后还想修改属性 `user.Set()`

或者 `user := UserModel.New().UserID()` 这样只能修改一个属性，怎么修改多个属性

```go
package UserModel

type UserModelImpl struct {
	UserID   int
	UserName string
}

// New 初始化实例
func New(attrs ...UserModelAttrFunc) *UserModelImpl {
	u := &UserModelImpl{}
	// 对 u 里每个属性进行初始化
	// 强制类型转化。
	UserModelAttrFuncs(attrs).Apply(u)
	return u
}

// Mutate 修改实例属性
func (this *UserModelImpl) Mutate(attrs ...UserModelAttrFunc) *UserModelImpl {
	UserModelAttrFuncs(attrs).Apply(this)
	return this
}
```

调用

```go
package main

import (
   "ginskill/src/models/UserModel"
   "github.com/gin-gonic/gin"
)

func main() {
   r := gin.New()
   r.GET("/", func(c *gin.Context) {
      user := UserModel.New().
         Mutate(UserModel.WithUserID(3)).
         Mutate(UserModel.WithUserName("custer"))
      c.JSON(200, user)
   })
   r.Run(":8080")
}
```

也可以写成这样

```go
user := UserModel.New().
   Mutate(UserModel.WithUserID(3),
      UserModel.WithUserName("custer"))
```

### 04. 参数验证和error处理:基本方法

修改模型 `model.go` 提交参数时表单的 `form` 是 `name`

```go
type UserModelImpl struct {
	UserID   int `json:"id"`
	UserName string `json:"name" form:"name" binding:"min=4"`
}
```

基本使用

```go
package main

import (
	"ginskill/src/models/UserModel"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.POST("/", func(c *gin.Context) {
		user := UserModel.New()
		if err := c.ShouldBind(user); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
		} else {
			c.JSON(200, user)
		}
	})
	r.Run(":8080")
}
```

#### 封装错误返回中间件

在 `src` 目录下新建 `common` 文件夹，并新建 `middlewares.go` 文件

```go
package common

import "github.com/gin-gonic/gin"

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				c.JSON(400, gin.H{"message": e})
			}
		}()
		c.Next()
	}
}
```

在 `main` 中使用中间件

```go
package main

import (
	"ginskill/src/common"
	"ginskill/src/models/UserModel"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(common.ErrorHandler())
	r.POST("/", func(c *gin.Context) {
		user := UserModel.New()
		if err := c.ShouldBind(user); err != nil {
			panic(err.Error())
		} else {
			c.JSON(200, user)
		}
	})
	r.Run(":8080")
}
```

#### 封装整个错误处理

新建文件夹 `result` , 新建文件`error_result.go`

```go
package result

type ErrorResult struct {
	err error
}

func (this *ErrorResult) Unwrap() interface{} {
	if this.err != nil {
		panic(this.err.Error())
	}
	return nil
}

func Result(err error) *ErrorResult {
	return &ErrorResult{err: err}
}
```

在 `main.go` 中就可以使用

```go
package main

import (
	"ginskill/src/common"
	"ginskill/src/models/UserModel"
	"ginskill/src/result"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(common.ErrorHandler())
	r.POST("/", func(c *gin.Context) {
		user := UserModel.New()
		result.Result(c.ShouldBind(user)).Unwrap()
		c.JSON(200, user)
	})
	r.Run(":8080")
}
```

