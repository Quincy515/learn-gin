package lib

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("token") == "" {
			c.AbortWithStatusJSON(400, gin.H{"message": "token required"})
		} else {
			c.Set("user_name", c.Request.Header.Get("token"))
			c.Next()
		}
	}
}

func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user_name")
		//access, err := E.Enforce(user, c.Request.RequestURI, c.Request.Method)
		domain := c.Param("domain")
		uri := strings.TrimPrefix(c.Request.RequestURI, "/"+domain) // /domain/depts => /depts
		access, err := E.Enforce(user, domain, uri, c.Request.Method)
		if err != nil || !access {
			c.AbortWithStatusJSON(403, gin.H{"message": "forbidden"})
		} else {
			c.Next()
		}
	}
}

func Middlewares() (fs []gin.HandlerFunc) {
	fs = append(fs, CheckLogin(), RBAC())
	return
}
