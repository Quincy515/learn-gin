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

// UserTest 控制器方法
func (this *UserClass) UserTest(c *gin.Context) string {
	return "用户测试"
}

// UserList 用户列表 返回切片
func (this *UserClass) UserList(c *gin.Context) goft.Models {
	users := []*models.UserModel{
		&models.UserModel{UserID: 101, UserName: "custer"},
		{UserID: 102, UserName: "张三"},
		{UserID: 103, UserName: "李四"},
	}
	return goft.MakeModels(users)
}

// UserDetail 返回 Model 实体(即所有自定义的 struct)，返回值都是 goft.Model
func (this *UserClass) UserDetail(c *gin.Context) goft.Model {
	user := models.NewUserModel()
	err := c.BindUri(user)
	goft.Error(err, "ID 参数不合法") // 如果出错会发生 panic，然后在中间件中处理返回 400
	return user
}

func (this *UserClass) Build(goft *goft.Goft) {
	goft.Handle("GET", "/test", this.UserTest).
		Handle("GET", "/userlist", this.UserList).
		Handle("GET", "/user/:id", this.UserDetail)
}
