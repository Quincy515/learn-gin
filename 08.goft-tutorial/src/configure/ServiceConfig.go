package configure

import (
	"goft-tutorial/src/daos"
	"goft-tutorial/src/service"
)

type ServiceConfig struct{}

func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{}
}

// 定义一个 Dao 对应一个 service
func (this *ServiceConfig) UserDao() *daos.UserDAO {
	return daos.NewUserDAO()
}

func (this *ServiceConfig) UserService() *service.UserService {
	return service.NewUserService()
}
