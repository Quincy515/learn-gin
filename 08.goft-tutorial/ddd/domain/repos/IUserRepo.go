package repos

import "goft-tutorial/ddd/domain/models"

// IUserRepo 用户相关的仓储定义
type IUserRepo interface {
	FindById(*models.UserModel) error
	FindByName(*models.UserModel) error
	SaveUser(*models.UserModel) error
	UpdateUser(*models.UserModel) error
	DeleteUser(*models.UserModel) error
}

// IUserLogRepo 日志相关的仓储定义
type IUserLogRepo interface {
	FindByName(name string) *models.UserLogModel
	SaveLog(model *models.UserLogModel) error
}
