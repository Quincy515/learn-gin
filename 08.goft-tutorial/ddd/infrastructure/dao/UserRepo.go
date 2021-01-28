package dao

import (
	"goft-tutorial/ddd/domain/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// FindByName 在这里实现具体业务操作
func (u *UserRepo) FindByName(name string) *models.UserModel {
	user := &models.UserModel{}
	if u.DB.Where("user_name=?", name).Find(user).Error != nil {
		return nil
	}
	return user
}
func (u *UserRepo) SaveUser(*models.UserModel) error   { return nil }
func (u *UserRepo) UpdateUser(*models.UserModel) error { return nil }
func (u *UserRepo) DeleteUser(*models.UserModel) error { return nil }
