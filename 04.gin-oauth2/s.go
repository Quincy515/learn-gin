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
	"time"
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
	// 根据授权码获取 token
	r.POST("/token", func(c *gin.Context) {
		err := srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			panic(err.Error())
		}
	})
	// 如果没有登录，则跳转登录界面
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
	// 根据 token 返回用户信息
	r.POST("/info", func(c *gin.Context) {
		token, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			panic(err.Error())
		}
		log.Println("token: ", token.GetAccess())
		ret := gin.H{
			"user_id": token.GetUserID(),
			"expire": int64(token.GetAccessCreateAt().
				Add(token.GetAccessExpiresIn()).
				Sub(time.Now()).Seconds())}
		c.JSON(200, ret)
	})
	// 注销 token
	r.GET("/logout", func(c *gin.Context) {
		token, err := srv.Manager.LoadAccessToken(c, c.Query("access_token"))
		if err != nil {
			panic(err.Error())
			return
		}
		utils.DeleteUserSession(c)
		err = manager.RemoveAccessToken(c, token.GetAccess())
		if err != nil {
			panic(err.Error())
			return
		}
		c.Redirect(301, c.Query("redirect_uri"))
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
