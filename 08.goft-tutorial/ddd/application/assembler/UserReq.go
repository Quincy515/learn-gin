package assembler

import (
	"github.com/go-playground/validator/v10"
	"goft-tutorial/ddd/application/dto"
	"goft-tutorial/ddd/domain/models"
)

type UserReq struct {
	v *validator.Validate
}

func (u *UserReq) D2M_UserModel(dto *dto.SimpleUserReq) *models.UserModel {
	err := u.v.Struct(dto)
	if err != nil {
		panic(err.Error())
	}
	return models.NewUserModel(models.WithUserID(dto.Id))
}
