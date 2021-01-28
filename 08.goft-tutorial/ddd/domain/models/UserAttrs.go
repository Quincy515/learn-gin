package models

import "goft-tutorial/ddd/domain/valueobjs"

type UserAttrFunc func(model *UserModel) // 设置 User 属性的 函数类型
type UserAttrFuncs []UserAttrFunc

// 传参数
func WithUserID(id int) UserAttrFunc {
	return func(u *UserModel) {
		u.UserID = id
	}
}

func WithUserName(name string) UserAttrFunc {
	return func(u *UserModel) {
		u.UserName = name
	}
}

func WithUserPass(pass string) UserAttrFunc {
	return func(u *UserModel) {
		u.UserPwd = pass
	}
}

func WithUserExtra(extra *valueobjs.UserExtra) UserAttrFunc {
	return func(u *UserModel) {
		u.Extra = extra
	}
}

// apply 方法 循环 UserAttrFuncs 内容执行函数
func (u UserAttrFuncs) apply(userModel *UserModel) {
	for _, f := range u {
		f(userModel)
	}
}
