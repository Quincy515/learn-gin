### 01. oauth2基本概念

OAuth2（开放授权）

> 是一个开放标准，允许用户授权第三方移动应用访问他们存储在另外的服务提供者上的信息，而不需要将用户名和密码提供给第三方移动应用或分享他们数据的所有内容，OAuth2.0是OAuth协议的延续版本，但不向后兼容OAuth 1.0即完全废止了OAuth1.0 
>
>  最常见的场景就是 QQ登录、微信登录、github登录等

https://github.com/go-oauth2/oauth2

`go get -u github.com/go-oauth2/oauth2/v4`

四种模式

**授权码模式** (Authorization Code) 一般只使用这个，后面3个基本不用~~~

简化模式 (Implicit)

密码模式 (Resource Owner Password Credentials)

客户端模式 (Client Credentials)

```shell
     +--------+                               +---------------+
     |        |--(A)- Authorization Request ->|   Resource    |
     |        |                               |     Owner     |
     |        |<-(B)-- Authorization Grant ---|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(C)-- Authorization Grant -->| Authorization |
     | Client |                               |     Server    |
     |        |<-(D)----- Access Token -------|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(E)----- Access Token ------>|    Resource   |
     |        |                               |     Server    |
     |        |<-(F)--- Protected Resource ---|               |
     +--------+                               +---------------+
```

<img src="../imgs/14_code.png" style="zoom:95%;" />

### 02. 客户端请求授权码：基本参数

请求参数

| 参数名称      | 参数含义                         | 是否必须 |
| ------------- | -------------------------------- | -------- |
| response_type | 授权类型，此处的值为code         | 必须     |
| client_id     | 客户端 ID                        | 必须     |
| redirect_uri  | 重定向 URI                       | 必须     |
| scope         | 申请的权限范围，多个逗号隔开     | 可选     |
| state         | 客户端的当前状态，可以指定任意值 | 可选     |

授权码返回

| 参数名称 | 参数含义                                          | 是否必须 |
| -------- | ------------------------------------------------- | -------- |
| code     | 授权码。认证服务器返回的授权码                    | 必须     |
| state    | 如果A中请求包含这个参数，资源服务器原封不动的返回 | 可选     |

客户端代码

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/url"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("public/*")
	codeUrl, _ := url.ParseRequestURI("http://localhost:8080/getcode")
	loginUrl := "http://127.0.0.1:8081/auth?" +
		"response_type=code&client_id=clienta&redirect_uri=" +
		codeUrl.String()

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "a-index.html", map[string]string{
			"loginUrl": loginUrl,
		})
	})
	r.GET("/getcode", func(c *gin.Context) {
		code, _ := c.GetQuery("code")
		c.JSON(200, gin.H{"code": code})
	})
	r.Run(":8080")
}
```

服务端代码

```dago
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"log"
	"net/http"
)

func main() {
	manager := manage.NewDefaultManager()                 // 1. 创建管理对象
	manager.MustTokenStorage(store.NewMemoryTokenStore()) // 保存 token

	clientStore := store.NewClientStore() // 2. 客户端和服务端关联
	err := clientStore.Set("clienta", &models.Client{
		ID:     "clienta", // a 网站 id
		Secret: "123",
		Domain: "http://localhost:8080",
	})
	if err != nil {
		log.Fatal(err)
	}
	manager.MapClientStorage(clientStore) // 映射 client store

	// 3. 创建 http server
	srv := server.NewDefaultServer(manager)
	srv.SetUserAuthorizationHandler(userAuthorizationHandler)

	r := gin.New()
	// 相应授权码
	r.GET("/auth", func(c *gin.Context) {
		err := srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			log.Println(err)
		}
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	r.LoadHTMLGlob("public/*.html")
	r.Run(":8081")
}

// userAuthorizationHandler 确保授权时跳转到登录页面 login.html
func userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	w.Header().Set("Location", "/login")
	w.WriteHeader(302)

	return
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/50fe1de9b43f3efd5777a0dd4331fbe9290d26e7#diff-da6e607e95bad42ffe69d9ebe05a342ad7ac388011ff6371da7515533a4d50f6R1)

### 03. 客户端请求授权码：获取授权码

上面实现了点击 `使用第三方账号登录` 跳转到第三方的服务，发现没有登录，因此跳转到登录页。

下面要实现的是在登录页输入正确的用户名和密码登录之后跳转到原先 a 网站的页面上。并显示出授权码。

一旦用户登录成功，使用 `session` 保存用户信息，这里使用第三方库 https://github.com/gorilla/sessions

`go get github.com/gorilla/sessions`

新建 `utils` 文件夹，新建 `session.go` 文件

```go
package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
)

var sessionStore = sessions.NewCookieStore([]byte("123456"))

func init() {
	sessionStore.Options.Domain = "127.0.0.1:8081"
	sessionStore.Options.Path = "/"
	sessionStore.Options.MaxAge = 0 // 关闭浏览器就清除 session
}

// SaveUserSession 保存当前用户的信息
func SaveUserSession(c *gin.Context, userID string) {
	s, err := sessionStore.Get(c.Request, "LoginUser")
	if err != nil {
		panic(err.Error())
	}
	s.Values["userID"] = userID
	err = s.Save(c.Request, c.Writer) // save 保存
	if err != nil {
		panic(err.Error())
	}
}

// GetUserSession 获取用户信息
func GetUserSession(r *http.Request) string {
	if s, err := sessionStore.Get(r, "LoginUser"); err == nil {
		if s.Values["userID"] != nil {
			return s.Values["userID"].(string)
		}
	}
	return ""
}
```

新建`middlewares.go`

```,zgo
package utils

import "github.com/gin-gonic/gin"

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				context.JSON(400, gin.H{"message": e})
			}
		}()
		context.Next()
	}
}
```

中间件的配置，在 `s.go` 中

```go
r.Use(utils.ErrorHandler())
```

修改登录的页面 `login.html`

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <style>
        .body{margin:50px auto;text-align: center}
    </style>
</head>
<body>
 <div class="body">
     <form method="post">
     <div>
         <span>用户名:</span>
         <input type="text" name="userName"/>
     </div>
     <div>
         <span>密码:</span>
         <input type="password" name="userPass"/>
     </div>
     <div><input type="submit" value="登录"/></div>
        <div style="color: red">{{.error}}</div>
     </form>
 </div>
</body>
</html>
```

然后把 

```go
r.GET("/login", func(c *gin.Context) {
   c.HTML(200, "login.html", nil)
})
```

修改为

```go
package main

import (
	"gin-oauth2/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"log"
	"net/http"
)

func main() {
	manager := manage.NewDefaultManager()                 // 1. 创建管理对象
	manager.MustTokenStorage(store.NewMemoryTokenStore()) // 保存 token

	clientStore := store.NewClientStore() // 2. 客户端和服务端关联
	err := clientStore.Set("clienta", &models.Client{
		ID:     "clienta", // a 网站 id
		Secret: "123",
		Domain: "http://localhost:8080",
	})
	if err != nil {
		log.Fatal(err)
	}
	manager.MapClientStorage(clientStore) // 映射 client store

	// 3. 创建 http server
	srv := server.NewDefaultServer(manager)
	srv.SetUserAuthorizationHandler(userAuthorizationHandler)

	r := gin.New()
	r.Use(utils.ErrorHandler())
	// 相应授权码
	r.GET("/auth", func(c *gin.Context) {
		err := srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			log.Println(err)
		}
	})
	r.Any("/login", func(c *gin.Context) {
		data := map[string]string{
			"error": "",
		}
		if c.Request.Method == http.MethodPost {
			// 对提交的信息进行处理
			uname, upass := c.PostForm("userName"), c.PostForm("userPass")
			if uname+upass == "custer123" {
				utils.SaveUserSession(c, uname)
				c.Redirect(302, "/auth?"+c.Request.URL.RawQuery)
				return
			} else {
				data["error"] = "用户名密码错误"
			}
		}
		c.HTML(200, "login.html", data)
	})
	r.LoadHTMLGlob("public/*.html")
	r.Run(":8081")
}

// userAuthorizationHandler 确保授权时跳转到登录页面 login.html
func userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if userID = utils.GetUserSession(r); userID == "" {
		w.Header().Set("Location", "/login?"+r.URL.RawQuery)
		w.WriteHeader(302)
	}
	return
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/bf67092b04c442ea0ef6288b41eaa52fe5310fba#diff-49dab0ade8de98d1e754c8aba53f725f0ebcb8a86e3a20ea8e901e9525f270d7L4)

### 04. 客户端使用授权码请求token

上面的代码登录之后跳转到 http://localhost:8080/getcode?code=MGY2NDZJZGYTMWIYNY0ZOGVLLWI1NJYTYMRJOTVKZTGXMZE0 并获取到了授权码 `code`

go 的 `oauth2` 的客户端库 https://github.com/golang/oauth2，在 `a` 网站里需要安装这个库，安装

`go get golang.org/x/oauth2`

客户端 `a` 的代码

```go
package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const authServerURL = "http://127.0.0.1:8081"

var (
	oauth2Config = oauth2.Config{
		ClientID:     "clienta",
		ClientSecret: "123",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/getcode",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/auth",  // 获取授权码地址
			TokenURL: authServerURL + "/token", // 获取 token 地址
		},
	}
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("public/*")
	//codeUrl, _ := url.ParseRequestURI("http://localhost:8080/getcode")
	//loginUrl := "http://127.0.0.1:8081/auth?" +
	//	"response_type=code&client_id=clienta&redirect_uri=" +
	//	codeUrl.String()
	// state 参数，传递给服务端，验证通过会原封不动的传回客户端，
	// 在 /getcode 里理论上需要对 state 进行判断，防止被串改。
	loginUrl := oauth2Config.AuthCodeURL("myclient")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "a-index.html", map[string]string{
			"loginUrl": loginUrl,
		})
	})
	r.GET("/getcode", func(c *gin.Context) {
		code, _ := c.GetQuery("code")
		//c.JSON(200, gin.H{"code": code})
		token, err := oauth2Config.Exchange(c, code)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
		} else {
			c.JSON(200, token)
		}
	})
	r.Run(":8080")
}
```

服务端 `s.go` 代码修改

```go
package main

import (
	"gin-oauth2/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"log"
	"net/http"
)

var srv *server.Server

func main() {
	manager := manage.NewDefaultManager()                 // 1. 创建管理对象
	manager.MustTokenStorage(store.NewMemoryTokenStore()) // 保存 token

	clientStore := store.NewClientStore() // 2. 客户端和服务端关联
	err := clientStore.Set("clienta", &models.Client{
		ID:     "clienta", // a 网站 id
		Secret: "123",
		Domain: "http://localhost:8080",
	})
	if err != nil {
		log.Fatal(err)
	}
	manager.MapClientStorage(clientStore) // 映射 client store

	// 3. 创建 http server
	srv = server.NewDefaultServer(manager)
	srv.SetUserAuthorizationHandler(userAuthorizationHandler)

	r := gin.New()
	r.Use(utils.ErrorHandler())
	// 相应授权码
	r.GET("/auth", func(c *gin.Context) {
		err := srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			log.Println(err)
		}
	})
	r.POST("/token", func(c *gin.Context) {
		err := srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			panic(err.Error())
		}
	})
	r.Any("/login", func(c *gin.Context) {
		data := map[string]string{
			"error": "",
		}
		if c.Request.Method == http.MethodPost {
			// 对提交的信息进行处理
			uname, upass := c.PostForm("userName"), c.PostForm("userPass")
			if uname+upass == "custer123" {
				utils.SaveUserSession(c, uname)
				c.Redirect(302, "/auth?"+c.Request.URL.RawQuery)
				return
			} else {
				data["error"] = "用户名密码错误"
			}
		}
		c.HTML(200, "login.html", data)
	})
	r.LoadHTMLGlob("public/*.html")
	r.Run(":8081")
}

// userAuthorizationHandler 确保授权时跳转到登录页面 login.html
func userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if userID = utils.GetUserSession(r); userID == "" {
		w.Header().Set("Location", "/login?"+r.URL.RawQuery)
		w.WriteHeader(302)
	}
	return
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/d2280d2bd5666b7fa5d14438f6a0b1457c7fd48d#diff-da6e607e95bad42ffe69d9ebe05a342ad7ac388011ff6371da7515533a4d50f6L2)

### 05. 客户端使用token获取用户信息

在服务端 `s.go` 中增加使用token获取用户信息的路由 `POST /info`

```go
	// 根据 token 返回用户信息
	r.POST("/info", func(c *gin.Context) {
		token, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			panic(err.Error())
		}
		ret := gin.H{
			"user_id": token.GetUserID(),
			"expire": int64(token.GetAccessCreateAt().
				Add(token.GetAccessExpiresIn()).
				Sub(time.Now()).Seconds())}
		c.JSON(200, ret)
	})
```

客户端请求用户信息的方式

OAuth 2.0 定义了获取 Access Token 的方法

第一种 :放在头里

```http
GET /info  HTTP/1.1
Host: server.example.com
Authorization: Bearer xxxxxooo
```

先通过 http://localhost:8080 点击第三方登录跳转到 服务端登录页面，登录成功后跳转回 http://localhost:8080/getcode?code=Y2QZMTFKYZYTZMQ1ZC0ZYJI0LTLJNGETZDM2NJZHMDJMNDC2&state=myclient，可以看到返回的 json

```json
{
    "access_token": "OTLIMGRLNWUTNWYYNC0ZZJLJLTHIZJQTMMUZMDQ1ZTLJYTZM",
    "token_type": "Bearer",
    "refresh_token": "NWRHYMQ2NDGTNWU4MC01YWJILWI0ZTATMJYZNGIWMZIXYZFM",
    "expiry": "2020-11-10T12:22:01.3254618+08:00"
}
```

发送 POST 请求

<img src="../imgs/14_token.png" style="zoom:95%;" />

第二种：放到Form参数里

```http
Content-Type: application/x-www-form-urlencoded
access_token=xxoxooxoo
```

<img src="../imgs/15_token.png" style="zoom:95%;" />

客户端可以请求`GET /info?access_token=`  修改代码  `a.go`

```go
	// 根据 token 返回用户信息
	r.POST("/info", func(c *gin.Context) {
		token, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			panic(err.Error())
		}
		ret := gin.H{
			"user_id": token.GetUserID(),
			"expire": int64(token.GetAccessCreateAt().
				Add(token.GetAccessExpiresIn()).
				Sub(time.Now()).Seconds())}
		c.JSON(200, ret)
	})
```

客户端写 `http.Requet` 请求服务端获取用户信息 `utils/user.go`

```go
package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func GetUserInfo(url string, token string, isBearer bool) string {
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)
	if isBearer {
		req, err = http.NewRequest("POST", url, nil)
		if err != nil {
			panic(err.Error())
		}
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		req, err = http.NewRequest("POST", url, strings.NewReader("access_token="+token))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			panic(err.Error())
		}
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	return string(b)
}
```

直接在浏览器访问 http://localhost:8080/info?token=MMVJYJU1ZWMTZGQ1YI0ZNZBILTHIOTMTYMFIZMIXOWUYNMJM 就可以查看到用户信息

```json
{
    "expire": 7032,
    "user_id": "custer"
}
```

代码变动 [git commit]()

