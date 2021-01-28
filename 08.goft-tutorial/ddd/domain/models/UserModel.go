package models

import (
	"crypto/md5"
	"fmt"
	"goft-tutorial/ddd/domain/valueobjs"
)

type UserModel struct {
	*Model
	UserID   int                  `gorm:"column:user_id" json:"user_id"`
	UserName string               `gorm:"column:user_name" json:"user_name"`
	UserPwd  string               `gorm:"column:user_pwd" json:"user_pwd"`
	Extra    *valueobjs.UserExtra // 值对象 - 通过属性指向用户的额外附加信息
}

// NewUserModel 构造函数
func NewUserModel(attrs ...UserAttrFunc) *UserModel {
	user := &UserModel{}
	UserAttrFuncs(attrs).apply(user)
	user.Model = &Model{}
	user.SetId(user.UserID)
	user.SetName("User Entity") // 用户实体名称
	return user
}

func (UserModel) TableName() string {
	return `user` //
}

func (u *UserModel) BeforeSave() {
	u.UserPwd = fmt.Sprintf("%x", md5.Sum([]byte(u.UserPwd)))
}
