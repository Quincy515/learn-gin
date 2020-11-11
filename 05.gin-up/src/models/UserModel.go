package models

import "fmt"

type UserModel struct {
	UserID   int `gorm:"column:user_id" uri:"id" binding:"required,gt=0"`
	UserName string
}

func (*UserModel) TableName() string {
	return "users"
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func (this *UserModel) String() string {
	return fmt.Sprintf("user_id: %d, user_name: %s", this.UserID, this.UserName)
}
