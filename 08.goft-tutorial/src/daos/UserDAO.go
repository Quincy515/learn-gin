package daos

import (
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/models"
)

type UserDAO struct {
	Db *XOrmAdapter `inject:"-"` // 依赖注入
}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

const getUserByID = `
			SELECT 
				user_id, user_name
			FROM 
				users
			WHERE
				user_id=?`

// 简单的查询使用返回 goft.Query, 以 Get 开头
func (this *UserDAO) GetUserByID(uid int) goft.Query {
	return goft.SimpleQuery(getUserByID).
		WithArgs(uid).WithFirst(). // WithArgs 返回包含对象的数组，WithFirst 直接返回第一个对象
		WithMapping(map[string]string{
			"user_id":   "userID",
			"user_name": "userName",
		})
}

// goft.Query 是给前端控制器使用的，一般不做为业务的控制
func (this *UserDAO) GetUserByName(uname string) goft.Query {
	return goft.SimpleQuery(`
			SELECT 
				user_id, user_name
			FROM 
				users
			WHERE
				user_name=?`).
		WithArgs(uname).WithFirst()
}

// orm 操作的函数都是以 FindBy 开头
func (this *UserDAO) FindByUserName(username string) *models.UserModel {
	userModel := &models.UserModel{}
	has, err := this.Db.Table("users").Where("user_name=?", username).Get(userModel)
	if err != nil || !has {
		panic("user not exists")
	}
	return userModel
}
