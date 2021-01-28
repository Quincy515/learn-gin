package assembler

import (
	"goft-tutorial/ddd/application/dto"
	"goft-tutorial/ddd/domain/aggregates"
	"goft-tutorial/ddd/domain/models"
)

// UserResp 把用户实体映射成 DTO 展示给前端
type UserResp struct{}

// M2D_SimpleUserInfo 把用户实体映射为 简单用户 DTO
func (u *UserResp) M2D_SimpleUserInfo(user *models.UserModel) *dto.SimpleUserInfo {
	simpleUser := &dto.SimpleUserInfo{}
	simpleUser.Id = user.UserID
	simpleUser.Name = user.UserName
	simpleUser.City = user.Extra.UserCity
	return simpleUser
}

func (u *UserResp) M2D_UserInfo(mem *aggregates.Member) *dto.UserInfo {
	userInfo := &dto.UserInfo{}
	userInfo.Id = mem.User.UserID
	///... 其他睡醒赋值
	userInfo.Logs = u.M2D_UserLogs(mem.GetLogs())
	return userInfo
}

func (u *UserResp) M2D_UserLogs(logs []*models.UserLogModel) (ret []*dto.UserLog) {
	return
}
