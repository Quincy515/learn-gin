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
