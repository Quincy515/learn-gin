package services

import (
	"fmt"
	"goft-tutorial/ddd/domain/repos"
	"goft-tutorial/ddd/infrastructure/utils"
)

type UserLoginService struct {
	userRepo repos.IUserRepo
}

// UserLogin 复杂的用户登录逻辑
func (u *UserLoginService) UserLogin(userName string, userPwd string) (string, error) {
	user := u.userRepo.FindByName(userName)
	if user.UserID > 0 { // 存在该用户
		if user.UserPwd == utils.Md5(userPwd) {
			// TODO：记录登录日志
			return "1000200", nil
		} else {
			return "1000400", fmt.Errorf("密码不正确")
		}
	} else {
		// 1000 代表用户，404代表不存在
		return "1000404", fmt.Errorf("用户不存在")
	}

}
