package services

import (
	"goft-tutorial/ddd/application/assembler"
	"goft-tutorial/ddd/application/dto"
	"goft-tutorial/ddd/domain/repos"
)

type UserService struct {
	assUserReq *assembler.UserReq
	assUserRsp *assembler.UserResp
	userRepo   repos.IUserRepo
}

func (u *UserService) GetSimpleUserInfo(req *dto.SimpleUserReq) *dto.SimpleUserInfo {
	userModel := u.assUserReq.D2M_UserModel(req)
	if err := u.userRepo.FindById(userModel); err != nil {
		return nil
	}
	return u.assUserRsp.M2D_SimpleUserInfo(userModel)
}
