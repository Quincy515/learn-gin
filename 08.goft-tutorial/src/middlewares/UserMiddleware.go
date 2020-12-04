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
