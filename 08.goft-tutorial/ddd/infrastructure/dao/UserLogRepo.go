package dao

import (
	"goft-tutorial/ddd/domain/models"
	"gorm.io/gorm"
)

type UserLogRepo struct {
	DB *gorm.DB
}

func (u *UserLogRepo) FindByName(name string) *models.UserLogModel { return nil }
func (u *UserLogRepo) SaveLog(model *models.UserLogModel) error    { return nil }
