package dao

import (
	"goft-tutorial/ddd/domain/models"
	"goft-tutorial/ddd/domain/repos"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

var _ repos.IUserRepo = &UserRepo{}

func (u *UserRepo) FindById(user *models.UserModel) error {
	return u.DB.Where("user_id=?", user.Id).Find(user).Error
}

// FindByName 在这里实现具体业务操作
func (u *UserRepo) FindByName(user *models.UserModel) error {
	return u.DB.Where("user_name=?", user.UserName).Find(user).Error
}
func (u *UserRepo) SaveUser(*models.UserModel) error   { return nil }
func (u *UserRepo) UpdateUser(*models.UserModel) error { return nil }
func (u *UserRepo) DeleteUser(*models.UserModel) error { return nil }
