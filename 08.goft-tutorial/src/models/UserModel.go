package models

// UserDetailRequest 用户请求实体 使用 gin 原生请求验证
type UserDetailRequest struct {
	UserId int `binding:"required,gt=0" uri:"uid"`
}

func NewUserDetailRequest() *UserDetailRequest {
	return &UserDetailRequest{}
}

type UserModel struct {
	UserId   int
	UserName string
}

func NewUserModel(userId int, userName string) *UserModel {
	return &UserModel{UserId: userId, UserName: userName}
}
