package handlers

import (
	"ginskill/src/data/getter"
	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	R(c)(
		"get_user_list",
		"10001",
		getter.UserGetter.GetUserList(),
	)(OK)
}
