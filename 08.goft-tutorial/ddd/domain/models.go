package models

import (
	"crypto/md5"
	"fmt"
)

type UserModel struct {
	UserID   int        `gorm:"column:user_id" json:"user_id"`
	UserName string     `gorm:"column:user_name" json:"user_name"`
	UserPwd  string     `gorm:"column:user_pwd" json:"user_pwd"`
	Extra    *UserExtra // 值对象 - 通过属性指向用户的额外附加信息
}
type UserExtra struct {
	UserPhone string `gorm:"column:user_phone" json:"user_phone"`
	UserCity  string `gorm:"column:user_city" json:"user_city"`
	UserQq    string `gorm:"column:user_qq" json:"user_qq"`
}

func (u UserExtra) Equals(other *UserExtra) bool {
	return u.UserPhone == other.UserPhone && u.UserQq == other.UserQq && u.UserCity == other.UserCity
}

func (UserModel) TableName() string {
	return `user` //
}

func (u *UserModel) BeforeSave() {
	u.UserPwd = fmt.Sprintf("%x", md5.Sum([]byte(u.UserPwd)))
}
