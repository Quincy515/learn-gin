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
