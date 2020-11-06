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

代码变动 [git commit]()

