package daos

import "goft-tutorial/pkg/goft"

type UserDAO struct{}

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

func (this *UserDAO) GetUserByID(uid int) goft.Query {
	return goft.SimpleQuery(getUserByID).
		WithArgs(uid).WithFirst(). // WithArgs 返回包含对象的数组，WithFirst 直接返回第一个对象
		WithMapping(map[string]string{
			"user_id":   "userID",
			"user_name": "userName",
		})
}

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
