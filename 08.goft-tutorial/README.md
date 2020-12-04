[toc]

jtthink 知识库 https://65480539.gitbook.io/jtthink/

Goft 脚手架使用文档 https://65480539.gitbook.io/goft/

Go微服务+领域驱动+K8s新版实训课开更(第一阶段)  https://65480539.gitbook.io/gop1/

### 01. 控制器的使用：返回String和JSON

基于成熟框架 `gin` 二次开发，或者在此上面做个脚手架 `goft`，定制业务等

<img src="../imgs/22_expr.jpg" style="zoom:90%;" />

安装 `go get -u github.com/shenyisyn/goft-gin@v0.4.1`

新建文件 `src/controllers/IndexController.go`

```go
package controllers

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
)

type IndexController struct{}

func NewIndexController() *IndexController {
	return &IndexController{}
}

// 返回 string
func (this *IndexController) Index(ctx *gin.Context) string {
	return "index"
}

// 返回 json
func (this *IndexController) IndexJSON(ctx *gin.Context) goft.Json {
  goft.Throw("测试异常", 500, ctx)
	return gin.H{"result": "index"}
}

func (this *IndexController) Name() string {
	return "IndexController"
}

func (this *IndexController) Build(goft *goft.Goft) {
	goft.Handle("GET", "/", this.IndexJSON)
}
```

新建启动程序 `main.go`

```go
package main

import (
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/controllers"
)

func main() {
	goft.Ignite().
		Mount("v1", controllers.NewIndexController()).
		Launch()
}
```

运行访问查看控制台

```bash
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /v1/                      --> goft-tutorial/pkg/goft.StringResponder.RespondTo.func1 (2 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```

访问页面 http://localhost:8080/v1/ 可以看到

- 返回的 `string ` 是`index` 
- 返回的 `json` 是 `{"result":"index"}`
- 返回的 `error` 是 `{"error":"测试异常"}`

代码初始化 [git commit](https://github.com/custer-go/learn-gin/commit/c1390e3fbe9a54c55c647494422474595b21f6e4#diff-5e031c8fe909e21e054d942a61a9503aad9eed28cc4d7bd5718110d4a74cd23eR1)

### 02. 中间件的使用(1)：判断必要参数

考虑可能整合其他的框架，所以没有使用 `gin` 自带的中间件。

所以 `goft` 实现了一个简易的中间件。

接口：

```go
type Fairing interface {
  OnRequest(*gin.Context) error
  OnResponse(result interface{})(interface{}, error)
}
```

`OnRequest` : 执行控制器方法前，修改如头信息、判断参数等等。

`OnResponse` : 执行控制器方法后，可以修改返回值内容。

只要实现了这两个方法一律视为中间件。

新建文件 `src/middlewares/TokenCheck.go`

```go
package middlewares

import "github.com/gin-gonic/gin"

type TokenCheck struct{}

func NewTokenCheck() *TokenCheck {
	return &TokenCheck{}
}

func (this *TokenCheck) OnRequest(ctx *gin.Context) error {}

func (this *TokenCheck) OnResponse(result interface{}) (interface{}, error) {}
```

新增判断是否登录的逻辑

```go
package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type TokenCheck struct{}

func NewTokenCheck() *TokenCheck {
	return &TokenCheck{}
}

func (this *TokenCheck) OnRequest(ctx *gin.Context) error {
	if ctx.Query("token") == "" {
		//goft.Throw("token required", 503, ctx) // 使用 throw 可以自定义 status code
		return fmt.Errorf("token required") // 自定义错误返回的 status code 是 400
	}
	return nil
}

func (this *TokenCheck) OnResponse(result interface{}) (interface{}, error) {
	return result, nil
}
```

把中间件增加到主函数 `main.go`

```go
package main

import (
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/controllers"
	"goft-tutorial/src/middlewares"
)

func main() {
	goft.Ignite().
		Attach(middlewares.NewTokenCheck()).
		Mount("v1", controllers.NewIndexController()).
		Launch()
}
```

访问页面 http://localhost:8080/v1/ 可以看到返回 400 `{"error":"token required"}`

访问页面 http://localhost:8080/v1/?token=123 可以看到正确的返回 `{"result":"index"}`

使用 `goft.Throw()` 自定义返回的 `status code` 可以看到页面返回 503

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/d7f1b95a28eee5674aa94a264c86182767950d9f#diff-5e031c8fe909e21e054d942a61a9503aad9eed28cc4d7bd5718110d4a74cd23eL3)

### 03. 中间件的使用(2)：修改响应内容

为了让返回的字段添加 `version`，可以使用 `OnResponse`，为了熟练使用中间件，

这里新建一个中间件 `src/middlewares/AddVersion.go`

```go
package middlewares

import (
	"github.com/gin-gonic/gin"
)

type AddVersion struct{}

func NewAddVersion() *AddVersion {
	return &AddVersion{}
}

func (this *AddVersion) OnRequest(ctx *gin.Context) error {
	return nil
}

func (this *AddVersion) OnResponse(result interface{}) (interface{}, error) {
	return result, nil
}
```

修改逻辑

```go
func (this *AddVersion) OnResponse(result interface{}) (interface{}, error) {
	if m, ok := result.(gin.H); ok {
		m["version"] = "0.4.1"
		return m, nil
	}
	return result, nil
}
```

在主函数 `main.go` 中增加中间件

```go
func main() {
	goft.Ignite().
		Attach(middlewares.NewTokenCheck(), middlewares.NewAddVersion()).
		Mount("v1", controllers.NewIndexController()).
		Launch()
}
```

访问 http://localhost:8080/v1/?token=123 可以看到通过中间件的方式**修改响应结果**

`{"result":"index","version":"0.4.1"}`

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/be5d2dd2702135835ca3e7c832d8739512cf6f09#diff-5e031c8fe909e21e054d942a61a9503aad9eed28cc4d7bd5718110d4a74cd23eL8)

### 04. 路由级的中间件(1):基本使用

`main.go` 中 `Attach` 是全局中间件，会在请求结束后，修改响应对象，加入 `version` 版本号。

```go
func main() {
	goft.Ignite().
		Attach(middlewares.NewTokenCheck(), middlewares.NewAddVersion()).
		Mount("v1", controllers.NewIndexController()).
		Launch()
}
```

下面针对 `/v1/test` 执行单独中间件，新增中间件 `src/middlewares/IndexTest.go`

```go
package middlewares

import "github.com/gin-gonic/gin"

type IndexTest struct{}

func NewIndexTest() *IndexTest {
	return &IndexTest{}
}

func (this *IndexTest) OnRequest(ctx *gin.Context) error {
	return nil
}

func (this *IndexTest) OnResponse(result interface{}) (interface{}, error) {
	if m, ok := result.(gin.H); ok {
		m["metadata"] = "index test"
		return m, nil
	}
	return result, nil
}
```

然后在控制器中修改

```go
package controllers

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/middlewares"
)

type IndexController struct{}

func NewIndexController() *IndexController {
	return &IndexController{}
}

// 返回 json
func (this *IndexController) Index(ctx *gin.Context) goft.Json {
	//goft.Throw("测试异常", 500, ctx)
	return gin.H{"result": "index"}
}

func (this *IndexController) IndexTest(ctx *gin.Context) goft.Json {
	return gin.H{"result": "index test"}
}

func (this *IndexController) Name() string {
	return "IndexController"
}

func (this *IndexController) Build(goft *goft.Goft) {
	goft.Handle("GET", "/", this.Index).
		HandleWithFairing("GET", "/test", this.IndexTest, middlewares.NewIndexTest())
}
```

针对路由的中间件拦截 `HandleWithFairing("GET", "/test", this.IndexTest, middlewares.NewIndexTest())`

访问页面 http://localhost:8080/v1/test?token=123 可以看到**路由级针对单独`url`中间件的执行 **

`{"metadata":"index test","result":"index test","version":"0.4.1"}`

代码修改 [git commit](https://github.com/custer-go/learn-gin/commit/8afdfce9d0ff57b6bb1393aa5e94b24c412be330#diff-18266a8616923f74411b54b15f0eb4eb72e8d9c6bfba34c4efeaf75aaa711d36L3)

### 05. 路由级的中间件(2):参数验证和业务分离（上）

场景 `GET /user/123` 得到用户 `ID = 123` 的用户信息。

新建一个用户控制器 `src/controllers/UserController.go`

```go
package controllers

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (this *UserController) Name() string {
	return "UserController"
}

func (this *UserController) Build(goft *goft.Goft) {
	goft.Handle("GET", "/user/:uid", this.UserDetail)
}

func (this *UserController) UserDetail(ctx *gin.Context) goft.Json {
}
```

新增用户的模型 `src/models/UserModel.go`

```go
package models

type UserModel struct {
	UserId   int
	UserName string
}
```

一般是在 `func (this *UserController) UserDetail(ctx *gin.Context) goft.Json {` 控制器中写参数验证。

相对正规的验证方法是先写请求实体，使用 `gin` 原生验证

```go
package models

// UserDetailRequest 用户请求实体 使用 gin 原生请求验证
type UserDetailRequest struct {
	UserId int `binding:"required,gt=0" uri:"uid"`
}

func NewUserDetailRequest() *UserDetailRequest {
	return &UserDetailRequest{}
}

type UserModel struct {
	UserId   int
	UserName string
}

func NewUserModel(userId int, userName string) *UserModel {
	return &UserModel{UserId: userId, UserName: userName}
}
```

然后写业务逻辑

```go
func (this *UserController) UserDetail(ctx *gin.Context) goft.Json {
	req := models.NewUserDetailRequest()
	goft.Error(ctx.BindUri(req)) // 出错就自动抛出异常，没有错误就继续执行
	return gin.H{"result": models.NewUserModel(req.UserId, "testUserName")}
}
```

把控制器增加到 `main.go` 中

```go
func main() {
	goft.Ignite().
		Attach(middlewares.NewTokenCheck(), middlewares.NewAddVersion()).
		Mount("v1", controllers.NewIndexController(),
			controllers.NewUserController()).
		Launch()
}
```

访问 http://localhost:8080/v1/user/123?token=1 可以看到结果

```json
{
  "result": {
    "UserId": 123,
    "UserName": "testUserName"
	},
	"version": "0.4.1"
}
```

下面分离代码在控制器中仅仅处理业务，验证部分代码可以专门封装到中间件中，一旦以后参数验证规则发生改变，就不需要更改 `controller` 代码，只需要修改中间件或者替换中间件。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/6b2e020b9ef9568c92c8c5dbb153ea467b1bcb13#diff-5e031c8fe909e21e054d942a61a9503aad9eed28cc4d7bd5718110d4a74cd23eL9)

### 06. 路由级的中间件(2):参数验证和业务分离（下）

由于现在有路由级的中间件，所以可以创建一个中间件 `src/middlewares/UserMiddleware.go`

```go
package middlewares

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/models"
)

type UserMiddleware struct{}

func NewUserMiddleware() *UserMiddleware {
	return &UserMiddleware{}
}

func (this *UserMiddleware) OnRequest(ctx *gin.Context) error {
	req := models.NewUserDetailRequest()
	goft.Error(ctx.BindUri(req))
	ctx.Set("_req", req)
	return nil
}

func (this *UserMiddleware) OnResponse(result interface{}) (interface{}, error) {
	return result, nil
}
```

参数验证放入中间件中，控制器只处理逻辑

```go
func (this *UserController) UserDetail(ctx *gin.Context) goft.Json {
	req, _ := ctx.Get("_req")
	return gin.H{"result": models.NewUserModel(req.(*models.UserDetailRequest).UserId, "testUserName")}
}
```

把中间件和路由匹配到一起

```go
func (this *UserController) Build(goft *goft.Goft) {
	goft.HandleWithFairing("GET", "/user/:uid", this.UserDetail, middlewares.NewUserMiddleware())
}
```

运行代码访问 http://localhost:8080/v1/user/2?token=1 可以看到和之前结果相同。

但是参数处理，和业务处理已经通过路由级中间件分离开来。

代码变动 [git commit]()