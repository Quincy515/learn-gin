package controllers

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/middlewares"
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
	goft.HandleWithFairing("GET", "/user/:uid", this.UserDetail, middlewares.NewUserMiddleware())
}

func (this *UserController) UserDetail(ctx *gin.Context) goft.Json {
	req, _ := ctx.Get("_req")
	return gin.H{"result": models.NewUserModel(req.(*models.UserDetailRequest).UserId, "testUserName")}
}
