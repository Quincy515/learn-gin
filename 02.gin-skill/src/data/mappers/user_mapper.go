package mappers

import (
	"ginskill/src/models/UserModel"
	"github.com/Masterminds/squirrel"
	"time"
)

type UserMapper struct{}

// GetUserList 获取用户列表
func (*UserMapper) GetUserList() *SqlMapper {
	return Mapper(squirrel.Select("user_id", "user_name").
		From("users").
		OrderBy("user_id desc").
		Limit(10).ToSql())
}

// GetUserDetail 获取用户详情
func (*UserMapper) GetUserDetail(id int) *SqlMapper {
	return Mapper(squirrel.Select("user_id", "user_name").
		From("users").Where("user_id=?", id).ToSql())
}

// AddNewUser 新增用户，传入用户实体
func (*UserMapper) AddNewUser(user *UserModel.UserModelImpl) *SqlMapper {
	return Mapper(squirrel.Insert(user.TableName()).
		Columns("user_name", "user_pwd", "user_addtime").
		Values(user.UserName, user.UserPwd, time.Now()).ToSql())
}
