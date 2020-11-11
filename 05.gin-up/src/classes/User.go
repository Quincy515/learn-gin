package classes

import (
	"gin-up/src/goft"
	"gin-up/src/models"
	"github.com/gin-gonic/gin"
)

// UserClass 
type UserClass struct{}

// NewUserClass UserClass generate constructor
func NewUserClass() *UserClass {
	return &UserClass{}
}

// UserList 控制器方法
func (this *UserClass) UserList(c *gin.Context) string {
	return "用户列表"
}

// UserDetail 返回 Model 实体
func (this *UserClass) UserDetail(c *gin.Context) goft.Model {
	return &models.UserModel{UserID: 101, UserName: "custer"}
}

func (this *UserClass) Build(goft *goft.Goft) {
	goft.Handle("GET", "/user", this.UserList).
		Handle("GET", "/user/detail", this.UserDetail)
}
