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

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/caa8cb0a29fd3407387a7dc7681568749db5d256#diff-fe3b020a336c7e0ea80e1ee4f700c33695d0a78c695d938e5b309e99e559e621L3)

### 07. 依赖注入和ORM 使用 (Gorm)

关于依赖注入，`goft` 使用的就是这个

**手撸IoC容器(golang)初级版本 **http://b.jtthink.com/read.php?tid=573&fid=2 

课程地址 https://www.jtthink.com/course/128

代码地址 https://github.com/shenyisyn/goft-ioc

连接字符串 --- 常规写法

新建文件夹 `src/configure/DBConfig.go`

```go
package configure

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type DBConfig struct{}

func NewDBConfig() *DBConfig {
	return &DBConfig{}
}

func (this *DBConfig) GormDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
			Colorful: true,        // 彩色打印
		},
	)
	dsn := "root:root1234@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB.SetMaxIdleConns(5)                   //最大空闲数
	mysqlDB.SetMaxOpenConns(10)                  //最大打开连接数
	mysqlDB.SetConnMaxLifetime(time.Second * 30) //空闲连接生命周期
	return db
}
```

修改 `main.go` 增加数据库连接的配置

```go
func main() {
	goft.Ignite().
		Config(configure.NewDBConfig()).
		Attach(middlewares.NewTokenCheck(), middlewares.NewAddVersion()).
		Mount("v1", controllers.NewIndexController(),
			controllers.NewUserController()).
		Launch()
}
```

修改`src/controller/UserController.go` 实现从数据库查询用户信息

```go
package controllers

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/middlewares"
	"goft-tutorial/src/models"
	"gorm.io/gorm"
)

type UserController struct {
	Db *gorm.DB `inject:"-"` // 依赖注入 - 表示单例模式
}

func NewUserController() *UserController {
	return &UserController{}
}

func (this *UserController) Name() string {
	return "UserController"
}

func (this *UserController) Build(goft *goft.Goft) {
	goft.HandleWithFairing("GET", "/user/:uid", this.UserDetail, middlewares.NewUserMiddleware())
}

func (this *UserController) UserDetail(ctx *gin.Context) goft.Json {
	req, _ := ctx.Get("_req")
	uid := req.(*models.UserDetailRequest).UserId
	user := &models.UserModel{}
	goft.Error(this.Db.Table("users").Where("user_id=?", uid).Find(user).Error)
	return user
}
```

访问 http://localhost:8080/v1/user/2?token=1 可以看到 `{"UserId":2,"UserName":"lisi"}`

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/4d8e0ff6954b8c450907907b0200d6d1f035c8df#diff-5e031c8fe909e21e054d942a61a9503aad9eed28cc4d7bd5718110d4a74cd23eL2)

### 08. ORM执行简化:直接返回SQL语句(GORM)

目录结构

```bash
├── README.md
├── go.mod
├── go.sum
├── main.go
└── src                            // 源码目录
    ├── configure                  // 若干个 config 对象
    │   └── DBConfig.go            // 返回需要注入到容器里的对象
    ├── controllers                // 控制器
    │   ├── IndexController.go
    ├── middlewares                // 中间件
    │   ├── AddVersion.go
    └── models                     // 模型包含请求实体和验证
        └── UserModel.go

```

上面访问 http://localhost:8080/v1/user/2?token=1 可以看到 `{"UserId":2,"UserName":"lisi"}`

修改代码`UserDetail ` `src/controllers/UserController.go` 直接返回 `SQL` 语句 `goft.SimpleQuery`

```go
func (this *UserController) UserDetail(ctx *gin.Context) goft.SimpleQuery {
	return "SELECT * FROM users WHERE user_id=2"
}
```

 代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/dee7d9d26daaab08cb76502bb7e1aa957f14d5d4#diff-fe3b020a336c7e0ea80e1ee4f700c33695d0a78c695d938e5b309e99e559e621L4)

### 09. ORM执行简化:控制器直接返回SQL语句(XORM)

上面在控制器中直接返回 `SQL` 语句，下面适配 `xorm`

首先安装 `xorm` `go get xorm.io/xorm` 在文件 `src/configure/DBConfig.go` 中，

初始化 `xorm`

```go
func (this *DBConfig) XOrm() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:root1234@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	engine.DB().SetMaxIdleConns(5)
	engine.DB().SetMaxOpenConns(10)
	return engine
}
```

修改 `UserController.go`

```go
func (this *UserController) UserDetail(ctx *gin.Context) goft.Json {
	user := &models.UserModel{}
	_, err := this.Db.Table("users").Where("user_id=?", 2).
		Get(user)
	goft.Error(err)
	return user
}
```

访问 http://localhost:8080/v1/user/2?token=1 可以看到 `{"user_id":2,"user_name":"lisi"}`

直接返回 `SQL` 语句

```go
func (this *UserController) UserDetail(ctx *gin.Context) goft.SimpleQuery {
	return "SELECT * FROM users WHERE user_id=2"
}
```

修改 `XormAdapter` 适配器

```go
type XOrmAdapter struct {
	*xorm.Engine
}

func (this *XOrmAdapter) DB() *sql.DB {
	return this.Engine.DB().DB
}

func (this *DBConfig) XOrm() *XOrmAdapter {
	engine, err := xorm.NewEngine("mysql", "root:root1234@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	engine.DB().SetMaxIdleConns(5)
	engine.DB().SetMaxOpenConns(10)
	return &XOrmAdapter{Engine: engine}
}
```

注入：

```go
type UserController struct {
	//Db *gorm.DB `inject:"-"` // 依赖注入 - 表示单例模式
	Db *configure.XOrmAdapter `inject:"-"`
}
```

这样就可以直接返回 `SQL` 语句或者 `JSON`

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/bb16f7e3d45fb284c29142b2c59cec7d094eb0c7#diff-2583763e3a634de0ff7874859a53ffbb3f375c2686c373cc1536ca5bae00b9b9L1)

### 10. 控制器返回SQL语句：支持参数

```go
func (this *UserController) UserDetail(ctx *gin.Context) goft.Query {
	return goft.SimpleQuery("SELECT * FROM users WHERE user_id=?").
		WithArgs(ctx.Param("uid")).WithFirst() // WithArgs 返回包含对象的数组，WithFirst 直接返回第一个对象
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/d3fbb67a9fa2518eac8d02364603847d3c75ec07#diff-fe3b020a336c7e0ea80e1ee4f700c33695d0a78c695d938e5b309e99e559e621L3)

### 11. 控制器返回SQL语句：支持自定义JSON字段

第一种方式 使用 `SQL`

```go
func (this *UserController) UserDetail(ctx *gin.Context) goft.Query {
	return goft.SimpleQuery(`
			SELECT 
				user_id as uid, user_name as uname
			FROM 
				users
			WHERE
				user_id=?`).
		WithArgs(ctx.Param("uid")).WithFirst() // WithArgs 返回包含对象的数组，WithFirst 直接返回第一个对象
}
```

访问 http://localhost:8080/v1/user/3?token=1 {"userID":"3","userName":"custer"}

第二种方式 使用 `WithMapping`

```go
func (this *UserController) UserDetail(ctx *gin.Context) goft.Query {
	return goft.SimpleQuery(`
			SELECT 
				user_id, user_name
			FROM 
				users
			WHERE
				user_id=?`).
		WithArgs(ctx.Param("uid")).WithFirst(). // WithArgs 返回包含对象的数组，WithFirst 直接返回第一个对象
		WithMapping(map[string]string{
			"user_id":   "userID",
			"user_name": "userName",
		})
}
```

访问http://localhost:8080/v1/user/3?token=1{"userID":"3","userName":"custer"}

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/b2ce5cd0cf23c19ee6a58f501790381bdec9d300#diff-fe3b020a336c7e0ea80e1ee4f700c33695d0a78c695d938e5b309e99e559e621L50)

### 12. DAO层示例：用户DAO的写法

`DAO` -- `data access object` 数据访问层，即写一个类 (struct), 把数据库操作的代码封装起来。

定位 -- 介于 `controller` 和 `service` 层之间。

新建文件 `src/daos/UserDAO.go`

```go
package daos

import "goft-tutorial/pkg/goft"

type UserDAO struct{}

func (this *UserDAO) GetUserDetail(uid interface{}) goft.Query {
	return goft.SimpleQuery(`
			SELECT 
				user_id, user_name
			FROM 
				users
			WHERE
				user_id=?`).
		WithArgs(uid).WithFirst(). // WithArgs 返回包含对象的数组，WithFirst 直接返回第一个对象
		WithMapping(map[string]string{
			"user_id":   "userID",
			"user_name": "userName",
		})
}
```

然后在 `src/controller/UserController.go` 中使用依赖注入

```go
package controllers

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/daos"
	"goft-tutorial/src/middlewares"
	"gorm.io/gorm"
)

type UserController struct {
	Db *gorm.DB `inject:"-"` // 依赖注入 - 表示单例模式
	//Db *configure.XOrmAdapter `inject:"-"`
	user *daos.UserDAO
}

func NewUserController() *UserController {
	return &UserController{}
}

func (this *UserController) Name() string {
	return "UserController"
}

func (this *UserController) Build(goft *goft.Goft) {
	goft.Handle("GET", "/users", this.UserList).
		HandleWithFairing("GET", "/user/:uid", this.UserDetail, middlewares.NewUserMiddleware())
}

func (this *UserController) UserList(ctx *gin.Context) goft.SimpleQuery {
	return "select * from users"
}

func (this *UserController) UserDetail(ctx *gin.Context) goft.Query {
	return this.user.GetUserDetail(ctx.Param("uid"))
}
```

这里控制器 `controller` 和 数据访问层 `dao` 直接关联，中间没有使用 `service` 层也是可以的。

在 `DAO` 层还要一种写法，把 `SQL` 语句写成常量单独放在一个文件中

```go
package daos

import "goft-tutorial/pkg/goft"

type UserDAO struct{}

const getUserByID = `
			SELECT 
				user_id, user_name
			FROM 
				users
			WHERE
				user_id=?`

func (this *UserDAO) GetUserDetail(uid interface{}) goft.Query {
	return goft.SimpleQuery(getUserByID).
		WithArgs(uid).WithFirst(). // WithArgs 返回包含对象的数组，WithFirst 直接返回第一个对象
		WithMapping(map[string]string{
			"user_id":   "userID",
			"user_name": "userName",
		})
}
```

访问 localhost:8080/v1/user/3?token=1 可以看到 `{"userID":"3","userName":"custer"}`,

如果访问 localhost:8080/v1/user/13?token=1 没有数据就返回的是 `[]`，

因此需要个 `service` 层专门进行处理判断，如果取不出怎么做，如果取出来了怎么做。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/66529bcd861b48f7aa33abccdced8741f885b0ce?branch=66529bcd861b48f7aa33abccdced8741f885b0ce&diff=split#diff-fe3b020a336c7e0ea80e1ee4f700c33695d0a78c695d938e5b309e99e559e621L3)

### 13. Service层示例：用户Service层的基本写法

`service` 层 -- 往往用来结合 `DAO` 处理业务相关操作。一般 `service` 只调用自己的 `DAO`。

如果想调用其他 `dao` ，则通过引入该层的 `service` 。

`Controller` 层负责定义路由和连接 `Service` 层的整合。

为了后面的依赖注入，在 `configuration` 中写 `ServiceConfig.go`，

专门把 `Bean` 注入到 `IOC` 容器。

```go
package configure

import (
	"goft-tutorial/src/daos"
	"goft-tutorial/src/service"
)

type ServiceConfig struct{}

func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{}
}

func (this *ServiceConfig) UserDao() *daos.UserDAO {
	return daos.NewUserDAO()
}

func (this *ServiceConfig) UserService() *service.UserService {
	return service.NewUserService()
}
```

因为要注入 `Bean` 中，所以要修改 `main.go`

```go
func main() {
	goft.Ignite().
		Config(configure.NewDBConfig(), configure.NewServiceConfig()).
		Attach(middlewares.NewTokenCheck(), middlewares.NewAddVersion()).
		Mount("v1", controllers.NewIndexController(),
			controllers.NewUserController()).
		Launch()
}
```

这样所有的 `service` `Bean` 都可以依赖注入。

在 `src` 中创建文件夹 `service` ，新建文件 `UserService.go`。

```go
package service

import (
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/daos"
	"strconv"
)

type UserService struct {
	UserDao *daos.UserDAO `inject:"-"`
}

func NewUserService() *UserService {
	return &UserService{}
}

func (this *UserService) GetUserDetail(param string) goft.Query {
	if uid, err := strconv.Atoi(param); err == nil {
		return this.UserDao.GetUserByID(uid)
	} else {
		return this.UserDao.GetUserByName(param)
	}
}
```

修改 `src/daos/UserDAO.go`

```go
package daos

import "goft-tutorial/pkg/goft"

type UserDAO struct{}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

const getUserByID = `
			SELECT 
				user_id, user_name
			FROM 
				users
			WHERE
				user_id=?`

func (this *UserDAO) GetUserByID(uid int) goft.Query {
	return goft.SimpleQuery(getUserByID).
		WithArgs(uid).WithFirst(). // WithArgs 返回包含对象的数组，WithFirst 直接返回第一个对象
		WithMapping(map[string]string{
			"user_id":   "userID",
			"user_name": "userName",
		})
}

func (this *UserDAO) GetUserByName(uname string) goft.Query {
	return goft.SimpleQuery(`
			SELECT 
				user_id, user_name
			FROM 
				users
			WHERE
				user_name=?`).
		WithArgs(uname).WithFirst()
}
```

修改 `UserController.go`

```go
package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/service"
	"gorm.io/gorm"
)

type UserController struct {
	Db *gorm.DB `inject:"-"` // 依赖注入 - 表示单例模式
	//Db *configure.XOrmAdapter `inject:"-"`
	UserService *service.UserService `inject:"-"`
}

func NewUserController() *UserController {
	return &UserController{}
}

func (this *UserController) Name() string {
	return "UserController"
}

func (this *UserController) Build(goft *goft.Goft) {
	goft.Handle("GET", "/users", this.UserList).
		//HandleWithFairing("GET", "/user/:uid", this.UserDetail, middlewares.NewUserMiddleware())
		Handle("GET", "/user/:uid", this.UserDetail)
}

func (this *UserController) UserList(ctx *gin.Context) goft.SimpleQuery {
	return "select * from users"
}

func (this *UserController) UserDetail(ctx *gin.Context) goft.Query {
	fmt.Println("uid:", ctx.Param("uid"))
	return this.UserService.GetUserDetail(ctx.Param("uid"))
}
```

访问 localhost:8080/v1/user/custer?token=1 看到 `{"user_id":"3","user_name":"custer"}`

访问 localhost:8080/v1/user/2?token=1 看到 `{"userID":"2","userName":"lisi"}`

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/ae545a8a44573b16145c3c70a602fdda5781fe54?branch=ae545a8a44573b16145c3c70a602fdda5781fe54&diff=split#diff-5e031c8fe909e21e054d942a61a9503aad9eed28cc4d7bd5718110d4a74cd23eL9)

### 14. Service层示例：用户登录示例

#### 第1步：Dao层 orm 注入

```go
package daos

import (
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/models"
)

type UserDAO struct {
	Db *XOrmAdapter `inject:"-"` // 依赖注入
}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

const getUserByID = `
			SELECT 
				user_id, user_name
			FROM 
				users
			WHERE
				user_id=?`

// 简单的查询使用返回 goft.Query, 以 Get 开头
func (this *UserDAO) GetUserByID(uid int) goft.Query {
	return goft.SimpleQuery(getUserByID).
		WithArgs(uid).WithFirst(). // WithArgs 返回包含对象的数组，WithFirst 直接返回第一个对象
		WithMapping(map[string]string{
			"user_id":   "userID",
			"user_name": "userName",
		})
}

// goft.Query 是给前端控制器使用的，一般不做为业务的控制
func (this *UserDAO) GetUserByName(uname string) goft.Query {
	return goft.SimpleQuery(`
			SELECT 
				user_id, user_name
			FROM 
				users
			WHERE
				user_name=?`).
		WithArgs(uname).WithFirst()
}

// orm 操作的函数都是以 findBy 开头
func (this *UserDAO) findByUserName(username string) *models.UserModel {
	userModel := &models.UserModel{}
	has, err := this.Db.Table("users").Where("user_name=?", username).Get(userModel)
	if err != nil || !has {
		panic("user not exists")
	}
	return userModel
}
```

#### 第2步 service 层

```go
func (this *UserService) UserLogin(uname string, uid int) string {
	if this.UserDao.FindByUserName(uname).UserId == uid {
		return "token" + uname
	}
	panic("error user access")
}
```

#### 第3步 controller 层

```go
func (this *UserController) Build(goft *goft.Goft) {
	goft.Handle("GET", "/users", this.UserList).
		//HandleWithFairing("GET", "/user/:uid", this.UserDetail, middlewares.NewUserMiddleware())
		Handle("GET", "/user/:uid", this.UserDetail).
		Handle("GET", "/access_token", this.UserAccessToken)
}

// UserAccessToken 获取用户登录 token / access_token?uname=888&uid=***
func (this *UserController) UserAccessToken(ctx *gin.Context) goft.Json {
	if uname, uid := ctx.Query("uname"), ctx.Query("uid"); uname != "" && uid != "" {
		id, _ := strconv.Atoi(uid)
		return gin.H{"token": this.UserService.UserLogin(uname, id)}
	}
	panic("error user access params")
}
```

访问 localhost:8080/v1/access_token?uname=custer&id=3 返回 `{"token":"tokencuster","version":"0.4.1"}`

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/e72e231b098f760d2b0b49c7e44bcc15eeed4474?branch=e72e231b098f760d2b0b49c7e44bcc15eeed4474&diff=split#diff-5e031c8fe909e21e054d942a61a9503aad9eed28cc4d7bd5718110d4a74cd23eL10)

### 15. ORM简化：自定义输出key、Query执行

特别简单的 `SQL` 语句，或者特别复杂的 `SQL` 语句，可以直接返回 `goft.SimpleQuery`

```go
func (this *UserController) UserList(ctx *gin.Context) goft.SimpleQuery {
	return "select * from users"
}
```

修改为 `map["result"]map[string]interface{}`

```go
func (this *UserController) UserList(ctx *gin.Context) goft.Json {
	//return "select * from users"
	return goft.SimpleQuery("select * from users").WithKey("result").Get()
}
```

访问 localhost:8080/v1/users 可以看到

```json
{
    "result": [
        {
            "tenant_id": "1",
            "user_id": "1",
            "user_name": "shenyi"
        },
        {
            "tenant_id": "2",
            "user_id": "2",
            "user_name": "lisi"
        },
        {
            "tenant_id": "2",
            "user_id": "3",
            "user_name": "custer"
        }
    ],
    "version": "0.4.1"
}
```

如果

```go
func (this *UserController) UserList(ctx *gin.Context) goft.Json {
	//return "select * from users"
	ret := goft.SimpleQuery("select * from users").WithKey("result").Get() // map["result"]map[string]interface{}
	return ret.(gin.H)["result"]
}
```

返回结果

```json
[
    {
        "tenant_id": "1",
        "user_id": "1",
        "user_name": "shenyi"
    },
    {
        "tenant_id": "2",
        "user_id": "2",
        "user_name": "lisi"
    },
    {
        "tenant_id": "2",
        "user_id": "3",
        "user_name": "custer"
    }
]
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/af47ba446a525bb17a08f81357acf12810ab7527?branch=af47ba446a525bb17a08f81357acf12810ab7527&diff=split#diff-fe3b020a336c7e0ea80e1ee4f700c33695d0a78c695d938e5b309e99e559e621L30)

### 16. 超简领域驱动模型入门(1)：基本分层

上面的业务实现是

```bash
Controller<--------------> Service层
（定义路由、Service层的整合）     ↑
                              |
                              |
                              ↓
                             DAO层
```

传统的三层`Controller` `Service` `DAO` 设计简单，后期业务复杂，可维护性比较低。

#### 领域驱动模型DDD

> 领域业务人员、产品或设计、程序员共同商定一个 **邻域模型**，
>
> 根据业务商定模型，使用通用语言进行描述，通过程序员实现具体代码。

#### 失血模型

```bash
领域模型(DM)         ------------->             业务对象(BO)

UserObject          ------------->          1. queryUserList
id int              ------------->          2. createUserList
name string         ------------->          3. findByUserID ...
```

#### 充血模型

```bash
领域模型(DM)                  业务层

UserObject                  比如 UserLogin调用了UserQuery和Update                                          
id int
name string
                --------->  
UserAdd()
UserDel()
UserQuery()
UserUpdate()
```

> **包括持久层的逻辑都定义在领域模型中**。业务层主要调用模型层完成业务的组合调用和事务的封装。

> **使用 DDD 一般利用充血模型来扩展不同的分层。**

#### 基本的分层(四层)

- `Infrastructure` 基础实施层
- `domain` 邻域层
- `application` 应用层
- `interfaces` 表示层，也叫用户界面层或接口层(接收用户请求、获取参数、进行判断Controller层)

```bash
Infrastructure -->  Interfaces
               -->  Application
               -->  Domain
```

#### `infrastructure` 基础实施层

与所有层进行交互

- 自己写的业务工具类 
- 配置信息
- 第三方库的集成和初始化
- 数据持久化机制等

所有层都可以调用基础实施层

#### `domain` 领域层

核心层，业务逻辑会在该层实现，比如包含

- 实体
- 值对象
- 聚合
- 工厂方法
- `Repository` 仓储实例

#### `Application` 应用层

连接 `domain` 和 `interfaces` 层，

对于 `interface` 层，提供各种业务功能方法给（Controller层即interfaces），

对于 `domain` 层，调用 `domain` 层完成任务逻辑。

> 对于业务代码简单使用传统三层比较合适，`controller` -> `service` -> `dao`

### 17. 领域层:用户实体编写和值对象(初步)

数据表

```sql
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) DEFAULT NULL,
  `user_pwd` varchar(255) DEFAULT NULL,
  `user_phone` varchar(255) DEFAULT NULL,
  `user_city` varchar(255) DEFAULT NULL,
  `user_qq` varchar(255) DEFAULT NULL,
  `user_addtime` datetime DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
```

#### `domain` 领域层

核心层，业务逻辑会在该层实现，比如包含

- 实体
- 值对象
- 聚合
- 工厂方法
- `Repository` 仓储实例

之前的代码都在 `src` ，现在在根目录下创建新的文件夹 `ddd/domain`

#### 实体

创建实体文件 `models.go` 包含用户实体。

```go
package models

import (
	"time"
)

type UserModel struct {
	UserID      int       `gorm:"column:user_id" json:"user_id"`
	UserName    string    `gorm:"column:user_name" json:"user_name"`
	UserPwd     string    `gorm:"column:user_pwd" json:"user_pwd"`
	UserPhone   string    `gorm:"column:user_phone" json:"user_phone"`
	UserCity    string    `gorm:"column:user_city" json:"user_city"`
	UserQq      string    `gorm:"column:user_qq" json:"user_qq"`
	UserAddtime time.Time `gorm:"column:user_addtime" json:"user_addtime"`
}

func (UserModel) TableName() string {
	return `user` //
}
```

实体第一要素，**要有唯一标识**

这里要做几件事

1. 定义用户实体 - 有唯一的标识（必须），包含各种属性，也可以包含如数据验证、操作前置函数，构造函数实例化等等
2. 用户的值对象

```go
// 前置操作、保存密码之前需要加密
func (u *UserModel) BeforeSave() { 
	u.UserPwd = fmt.Sprintf("%x", md5.Sum([]byte(u.UserPwd)))
}
```

#### 值对象

```go
type UserModel struct {
	UserID   int        `gorm:"column:user_id" json:"user_id"`
	UserName string     `gorm:"column:user_name" json:"user_name"`
	UserPwd  string     `gorm:"column:user_pwd" json:"user_pwd"`
	Extra    *UserExtra // 值对象 - 通过属性指向用户的额外附加信息
}
type UserExtra struct {
	UserPhone string `gorm:"column:user_phone" json:"user_phone"`
	UserCity  string `gorm:"column:user_city" json:"user_city"`
	UserQq    string `gorm:"column:user_qq" json:"user_qq"`
}
```

> 用来描述一个事物的特征，没有唯一标识的对象，譬如用户的 extra 信息

有2个重要原则

1. 实体只能通过 `ID`(唯一标识)来判断两者是否相同
2. 而值对象。只需根据“值”就能判断两者是否相等

不可变：修改值对象，必须传入新对象。

```go
func (u UserExtra) Equals(other *UserExtra) bool {
	return u.UserPhone == other.UserPhone && u.UserQq == other.UserQq && u.UserCity == other.UserCity
}
```

全部 `models.go`代码

```go
package models

import (
	"crypto/md5"
	"fmt"
)

type UserModel struct {
	UserID   int        `gorm:"column:user_id" json:"user_id"`
	UserName string     `gorm:"column:user_name" json:"user_name"`
	UserPwd  string     `gorm:"column:user_pwd" json:"user_pwd"`
	Extra    *UserExtra // 值对象 - 通过属性指向用户的额外附加信息
}
type UserExtra struct {
	UserPhone string `gorm:"column:user_phone" json:"user_phone"`
	UserCity  string `gorm:"column:user_city" json:"user_city"`
	UserQq    string `gorm:"column:user_qq" json:"user_qq"`
}

func (u UserExtra) Equals(other *UserExtra) bool {
	return u.UserPhone == other.UserPhone && u.UserQq == other.UserQq && u.UserCity == other.UserCity
}

func (UserModel) TableName() string {
	return `user` //
}

func (u *UserModel) BeforeSave() {
	u.UserPwd = fmt.Sprintf("%x", md5.Sum([]byte(u.UserPwd)))
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/2ee88ba0d2189cba6625f99f257f77f78db55b8c#diff-d0ccff768dc2a505a59253d2196158bd03e8303e5437a2bde670fe6dbff39d47R1)

### 18. 领域层:用户实体和值对象（2）--构造函数

上面创建了用户实体和值对象。上面把所有代码都放在了 `domain/models.go` 文件中。

这里修改下目录结构，新建文件夹 `domain/valueobjs` 和 `domain/models`。

```go
└── domain
		├── models
    │   └── UserModel.go
    └── valueobjs
        └── UserExtra.go
```

实体 `UserModel.go` 文件

```go
package models

import (
	"crypto/md5"
	"fmt"
	"goft-tutorial/ddd/domain/valueobjs"
)

type UserModel struct {
	UserID   int                  `gorm:"column:user_id" json:"user_id"`
	UserName string               `gorm:"column:user_name" json:"user_name"`
	UserPwd  string               `gorm:"column:user_pwd" json:"user_pwd"`
	Extra    *valueobjs.UserExtra // 值对象 - 通过属性指向用户的额外附加信息
}

func (UserModel) TableName() string {
	return `user` //
}

func (u *UserModel) BeforeSave() {
	u.UserPwd = fmt.Sprintf("%x", md5.Sum([]byte(u.UserPwd)))
}
```

值对象 `UserExtra.go` 文件

```go
package valueobjs

type UserExtra struct {
	UserPhone string `gorm:"column:user_phone" json:"user_phone"`
	UserCity  string `gorm:"column:user_city" json:"user_city"`
	UserQq    string `gorm:"column:user_qq" json:"user_qq"`
}

func (u UserExtra) Equals(other *UserExtra) bool {
	return u.UserPhone == other.UserPhone && u.UserQq == other.UserQq && u.UserCity == other.UserCity
}
```

创建构造函数文件   `domain/models/UserAttrs.go`

```go
package models

type UserAttrFunc func(model *UserModel) // 设置 User 属性的 函数类型
type UserAttrFuncs []UserAttrFunc

// 传参数
func WithUserID(id int) UserAttrFunc {
	return func(u *UserModel) {
		u.UserID = id
	}
}

func WithUserName(name string) UserAttrFunc {
	return func(u *UserModel) {
		u.UserName = name
	}
}

func WithUserPass(pass string) UserAttrFunc {
	return func(u *UserModel) {
		u.UserPwd = pass
	}
}

// apply 方法 循环 UserAttrFuncs 内容执行函数
func (u UserAttrFuncs) apply(userModel *UserModel) {
	for _, f := range u {
		f(userModel)
	}
}
```

在 `domain/models/UserModel.go` 中实现构造函数

```go
// NewUserModel 构造函数
func NewUserModel(attrs ...UserAttrFunc) *UserModel {
	user := &UserModel{}
	UserAttrFuncs(attrs).apply(user)
	return user
}
```

测试函数

```go
func main() {
  user := models.NewUserModel(
    models.WithUserName("custer"),
  )
  fmt.Println(user)
}
```

这就是在邻域模型里面创建构造函数的推荐方式。

同理在值对象中也新建构造函数的文件 `domain/valueobjs/UserExtraAttr.go`

```go
package valueobjs

type UserExtraAttrFunc func(model *UserExtra) // 设置 User 属性的 函数类型
type UserExtraAttrFuncs []UserExtraAttrFunc

// 传参数
func WithUserPhone(phone string) UserExtraAttrFunc {
	return func(u *UserExtra) {
		u.UserPhone = phone
	}
}

func WithUserQQ(qq string) UserExtraAttrFunc {
	return func(u *UserExtra) {
		u.UserQq = qq
	}
}

func WithUserCity(city string) UserExtraAttrFunc {
	return func(u *UserExtra) {
		u.UserCity = city
	}
}

// apply 方法 循环 UserExtraAttrFuncs 内容执行函数
func (u UserExtraAttrFuncs) apply(model *UserExtra) {
	for _, f := range u {
		f(model)
	}
}
```

在文件 `domain/valueobjs/UserExtra.go` 中实现构造函数

```go
func NewUserExtra(attrs ...UserExtraAttrFunc) *UserExtra {
	extra := &UserExtra{}
	UserExtraAttrFuncs(attrs).apply(extra)
	return extra
}
```

因为值对象不能直接传递，所以需要在 `domain/models/UserAttrs.go` 中修改代码

```go
func WithUserExtra(extra *valueobjs.UserExtra) UserAttrFunc {
	return func(u *UserModel) {
		u.Extra = extra
	}
}
```

测试函数

```go
func main() {
  user := models.NewUserModel(
    models.WithUserName("custer"),
    models.WithUserExtra(
      valueobjs.NewUserExtra(
        valueobjs.WithUserQQ("qqq"),
				valueobjs.WithUserQQ("上海"),
      )
    ),
  )
  fmt.Println(user)
	fmt.Println(user.Extra)
}
```

代码变动 [git commit]()