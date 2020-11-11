package goft

import "github.com/gin-gonic/gin"

// ErrorHandler 中间件
func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				context.AbortWithStatusJSON(400, gin.H{"error": e})
			}
		}()
		context.Next()
	}
}

// Error 出错直接 panic 然后在中间件中拦截
func Error(err error, msg ...string) {
	if err == nil {
		return
	} else {
		errMsg := err.Error() // 默认为内部的错误信息
		if len(msg) > 0 {     // 有自定义的报错信息
			errMsg = msg[0]
		}
		panic(errMsg)
	}
}
