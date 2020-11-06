package handlers

import (
	"ginskill/src/models/UserModel"
	"ginskill/src/result"
	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	user := UserModel.New()
	result.Result(c.ShouldBind(user)).Unwrap()
	if user.UserID > 10 {
		R(c)("user_list", "10000", "user_list")(OK)
	} else {
		R(c)("user_list", "10000", "error")(Error)
	}
}
