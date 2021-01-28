package aggregates

import (
	"goft-tutorial/ddd/domain/models"
	"goft-tutorial/ddd/domain/repos"
)

// Member 会员聚合 -- 会员：用户+日志+...组成
type Member struct {
	User *models.UserModel
	Log  *models.UserLogModel
	// 充值、社交、隐私信息
	userRepo    repos.IUserRepo // 接口
	userLogRepo repos.IUserLogRepo
}

// NewMember 构造函数
func NewMember(user *models.UserModel, userRepo repos.IUserRepo, userLogRepo repos.IUserLogRepo) *Member {
	return &Member{User: user, userRepo: userRepo, userLogRepo: userLogRepo}
}

// NewMemberByName 用户名作为唯一标识的构造函数
func NewMemberByName(name string, userRepo repos.IUserRepo, userLogRepo repos.IUserLogRepo) *Member {
	user := userRepo.FindByName(name)
	return &Member{User: user, userRepo: userRepo, userLogRepo: userLogRepo}
}

// Create 创建会员
func (m *Member) Create() error {
	err := m.userRepo.SaveUser(m.User)
	if err != nil {
		return err
	}
	m.Log = models.NewUserLogModel(m.User.UserName,
		models.WithUserLogType(models.UserLog_Create),
		models.WithUserLogComment("新增用户会员: "+m.User.UserName))
	return m.userLogRepo.SaveLog(m.Log)
}
