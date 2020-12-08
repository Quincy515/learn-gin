package controllers

import (
	"github.com/gin-gonic/gin"
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/middlewares"
	"gorm.io/gorm"
)

type UserController struct {
	Db *gorm.DB `inject:"-"` // 依赖注入 - 表示单例模式
	//Db *configure.XOrmAdapter `inject:"-"`
}

func NewUserController() *UserController {
	return &UserController{}
}

func (this *UserController) Name() string {
	return "UserController"
}

func (this *UserController) Build(goft *goft.Goft) {
	goft.Handle("GET", "/users", this.UserList).
		HandleWithFairing("GET", "/user/:uid", this.UserDetail, middlewares.NewUserMiddleware())
}

//func (this *UserController) UserDetail(ctx *gin.Context) goft.Json {
//	req, _ := ctx.Get("_req")
//	uid := req.(*models.UserDetailRequest).UserId
//	user := &models.UserModel{}
//	goft.Error(this.Db.Table("users").Where("user_id=?", uid).Find(user).Error)
//	return user
//}

//func (this *UserController) UserDetail(ctx *gin.Context) goft.SimpleQuery {
//	return goft.SimpleQuery(fmt.Sprintf("SELECT * FROM users WHERE user_id=%d", 3))
//}

//func (this *UserController) UserDetail(ctx *gin.Context) goft.Json {
//	user := &models.UserModel{}
//	_, err := this.Db.Table("users").Where("user_id=?", 2).
//		Get(user)
//	goft.Error(err)
//	return user
//}

func (this *UserController) UserList(ctx *gin.Context) goft.SimpleQuery {
	return "select * from users"
}

func (this *UserController) UserDetail(ctx *gin.Context) goft.Query {
	return goft.SimpleQuery(`
			SELECT 
				user_id, user_name
			FROM 
				users
			WHERE
				user_id=?`).
		WithArgs(ctx.Param("uid")).WithFirst(). // WithArgs 返回包含对象的数组，WithFirst 直接返回第一个对象
		WithMapping(map[string]string{
			"user_id":   "userID",
			"user_name": "userName",
		})
}
