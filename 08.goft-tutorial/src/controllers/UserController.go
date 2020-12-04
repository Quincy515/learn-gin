package controllers

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/models"
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
	req := models.NewUserDetailRequest()
	goft.Error(ctx.BindUri(req)) // 出错就自动抛出异常，没有错误就继续执行
	return gin.H{"result": models.NewUserModel(req.UserId, "testUserName")}
}
