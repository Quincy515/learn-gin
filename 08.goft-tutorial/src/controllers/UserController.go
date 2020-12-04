package controllers

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/middlewares"
	"goft-tutorial/src/models"
	"gorm.io/gorm"
)

type UserController struct {
	Db *gorm.DB `inject:"-"` // 依赖注入 - 表示单例模式
}

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
	uid := req.(*models.UserDetailRequest).UserId
	user := &models.UserModel{}
	goft.Error(this.Db.Table("users").Where("user_id=?", uid).Find(user).Error)
	return user
}
