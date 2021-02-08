package services

import (
	"goft-tutorial/ddd/application/assembler"
	"goft-tutorial/ddd/application/dto"
	"goft-tutorial/ddd/domain/repos"
)

type UserService struct {
	AssUserReq *assembler.UserReq
	AssUserRsp *assembler.UserResp
	userRepo   repos.IUserRepo `inject:"-"`
}

func (u *UserService) GetSimpleUserInfo(req *dto.SimpleUserReq) *dto.SimpleUserInfo {
	userModel := u.AssUserReq.D2M_UserModel(req)
	if err := u.userRepo.FindById(userModel); err != nil {
		return nil
	}
	//userModel.UserID = userModel.Id
	//userModel.UserName = "custer"
	//userModel.UserPwd = "123"
	//userModel.Extra.UserCity = "上海"
	return u.AssUserRsp.M2D_SimpleUserInfo(userModel)
}
