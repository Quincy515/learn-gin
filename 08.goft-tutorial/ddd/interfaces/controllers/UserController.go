package controllers

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/ddd/application/dto"
	"goft-tutorial/ddd/application/services"
	"goft-tutorial/ddd/infrastructure/utils"
	"goft-tutorial/pkg/goft"
)

type UserController struct {
	UserSvr *services.UserService `inject:"-"`
}

func NewUserController() *UserController {
	return &UserController{}
}

// UserDetail GET /user/123
func (u *UserController) UserDetail(ctx *gin.Context) goft.Json {
	//simpleUserReq := &dto.SimpleUserReq{}
	//ctx.ShouldBindUri(simpleUserReq)
	//return u.UserSvr.GetSimpleUserInfo(simpleUserReq)

	return u.UserSvr.GetSimpleUserInfo(
		utils.Exec(ctx.ShouldBindUri, &dto.SimpleUserReq{}).
			Unwrap().(*dto.SimpleUserReq))
}

func (u *UserController) Build(goft *goft.Goft) {
	goft.Handle("GET", "/users/:id", u.UserDetail)
}

func (u *UserController) Name() string {
	return "UserController"
}
