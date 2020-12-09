package service

import (
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/daos"
	"strconv"
)

type UserService struct {
	UserDao *daos.UserDAO `inject:"-"`
}

func NewUserService() *UserService {
	return &UserService{}
}

func (this *UserService) GetUserDetail(param string) goft.Query {
	if uid, err := strconv.Atoi(param); err == nil {
		return this.UserDao.GetUserByID(uid)
	} else {
		return this.UserDao.GetUserByName(param)
	}
}
