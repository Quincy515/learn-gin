package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/service"
	"gorm.io/gorm"
)

type UserController struct {
	Db *gorm.DB `inject:"-"` // 依赖注入 - 表示单例模式
	//Db *configure.XOrmAdapter `inject:"-"`
	UserService *service.UserService `inject:"-"`
}

func NewUserController() *UserController {
	return &UserController{}
}

func (this *UserController) Name() string {
	return "UserController"
}

func (this *UserController) Build(goft *goft.Goft) {
	goft.Handle("GET", "/users", this.UserList).
		//HandleWithFairing("GET", "/user/:uid", this.UserDetail, middlewares.NewUserMiddleware())
		Handle("GET", "/user/:uid", this.UserDetail)
}

func (this *UserController) UserList(ctx *gin.Context) goft.SimpleQuery {
	return "select * from users"
}

func (this *UserController) UserDetail(ctx *gin.Context) goft.Query {
	fmt.Println("uid:", ctx.Param("uid"))
	return this.UserService.GetUserDetail(ctx.Param("uid"))
}
