package handlers

import (
	"ginskill/src/models/UserModel"
	"ginskill/src/result"
	"ginskill/src/test"
	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	user := UserModel.New()
	result.Result(c.ShouldBind(user)).Unwrap()
	OK(c)("user_list", "10000", result.Result(test.GetInfo(user.UserID)).Unwrap())
}
