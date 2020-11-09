package setter

import (
	"ginskill/src/data/mappers"
	"ginskill/src/models/UserModel"
	"ginskill/src/result"
)

var UserSetter IUserSetter

func init() {
	UserSetter = NewUserSetterImpl()
}

type IUserSetter interface {
	SaveUser(*UserModel.UserModelImpl) *result.ErrorResult
}

type UserSetterImpl struct {
	userMapper *mappers.UserMapper
}

func NewUserSetterImpl() *UserSetterImpl {
	return &UserSetterImpl{userMapper: &mappers.UserMapper{}}
}

func (this *UserSetterImpl) SaveUser(user *UserModel.UserModelImpl) *result.ErrorResult {
	ret := this.userMapper.AddNewUser(user).Exec()
	return result.Result(ret.RowsAffected, ret.Error)
}
