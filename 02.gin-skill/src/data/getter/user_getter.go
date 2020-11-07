package getter

import (
	"fmt"
	"ginskill/src/data/mappers"
	"ginskill/src/models/UserModel"
	"ginskill/src/result"
)

// 对外使用的接口
var UserGetter IUserGetter

func init() {
	UserGetter = NewUserGetterImpl() // 业务更改，可以更换实现类
}

// IUserGetter 接口
type IUserGetter interface {
	GetUserList() []*UserModel.UserModelImpl // 返回实体列表
	GetUserByID(id int) *result.ErrorResult
}

// UserGetterImpl 实现 IUserGetter 接口
type UserGetterImpl struct {
	// 注入 UserMapper
	userMapper *mappers.UserMapper
}

// NewUserGetterImpl IUserGetter 接口的实现类
func NewUserGetterImpl() *UserGetterImpl {
	return &UserGetterImpl{userMapper: &mappers.UserMapper{}} // 在构造函数中赋值
}

// GetUserList 实现
func (this *UserGetterImpl) GetUserList() (users []*UserModel.UserModelImpl) {
	//dbs.Orm.Find(&users)
	//sqlMapper := this.userMapper.GetUserList()
	//dbs.Orm.Raw(sqlMapper.Sql, sqlMapper.Args).Find(&users) // 不管这个 ORM 是 gorm 还是 xorm 都可以执行原生 sql
	this.userMapper.GetUserList().Query().Find(&users)
	return
}

// GetUserByID 通过 id 获取 user 数据
func (this *UserGetterImpl) GetUserByID(id int) *result.ErrorResult {
	user := UserModel.New()
	//db := dbs.Orm.Where("user_id=?", id).Find(user)
	db := this.userMapper.GetUserDetail(id).Query().Find(user)
	if db.Error != nil || db.RowsAffected == 0 {
		return result.Result(nil, fmt.Errorf("not found user, id = %d", id))
	}
	return result.Result(user, nil)
}
