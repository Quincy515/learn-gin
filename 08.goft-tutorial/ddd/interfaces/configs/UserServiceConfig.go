package configs

import (
	"goft-tutorial/ddd/application/assembler"
	"goft-tutorial/ddd/application/services"
)

type UserServiceConfig struct{}

func NewUserServiceConfig() *UserServiceConfig {
	return &UserServiceConfig{}
}

// UserService 只会和 application 层相关，不会穿透到 domain 层
func (u *UserServiceConfig) UserService() *services.UserService {
	return &services.UserService{
		AssUserReq: &assembler.UserReq{},
		AssUserRsp: &assembler.UserResp{},
	}
}
