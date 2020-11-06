package getter

import (
	"ginskill/src/dbs"
	"ginskill/src/models/UserModel"
)

// 对外使用的接口
var UserGetter IUserGetter

func init() {
	UserGetter = NewUserGetterImpl() // 业务更改，可以更换实现类
}

// IUserGetter 接口
type IUserGetter interface {
	GetUserList() []*UserModel.UserModelImpl // 返回实体列表
}

// UserGetterImpl 实现 IUserGetter 接口
type UserGetterImpl struct{}

// NewUserGetterImpl IUserGetter 接口的实现类
func NewUserGetterImpl() *UserGetterImpl {
	return &UserGetterImpl{}
}

// GetUserList 实现
func (this *UserGetterImpl) GetUserList() (users []*UserModel.UserModelImpl) {
	dbs.Orm.Find(&users)
	return
}
