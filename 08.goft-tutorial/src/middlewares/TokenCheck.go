package middlewares

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
)

type TokenCheck struct{}

func NewTokenCheck() *TokenCheck {
	return &TokenCheck{}
}

func (this *TokenCheck) OnRequest(ctx *gin.Context) error {
	if ctx.Query("token") == "" {
		goft.Throw("token required", 503, ctx) // 使用 throw 可以自定义 status code
		//return fmt.Errorf("token required") // 自定义错误返回的 status code 是 400
	}
	return nil
}

func (this *TokenCheck) OnResponse(result interface{}) (interface{}, error) {
	return result, nil
}
