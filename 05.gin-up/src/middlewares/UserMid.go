package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// UserMid 用户中间件，"继承" Fairing 接口
type UserMid struct{}

func NewUserMid() *UserMid {
	return &UserMid{}
}

// OnRequest 在请求进入时，可以处理一些业务逻辑，或控制
func (this *UserMid) OnRequest(ctx *gin.Context) error {
	fmt.Println("这是新的用户中间件")
	fmt.Println(ctx.Query("name")) // 测试 query 参数，如 /v1/user?name=
	return nil
}
