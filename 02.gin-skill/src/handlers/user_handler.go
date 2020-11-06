package handlers

import (
	"ginskill/src/data/getter"
	"ginskill/src/result"
	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	R(c)(
		"get_user_list",
		"10001",
		getter.UserGetter.GetUserList(),
	)(OK)
}

// UserDetail 获取用户详情
func UserDetail(c *gin.Context) {
	// id := c.Param("id") 判断
	// 1. 获取 id
	id := &struct { // 使用匿名 struct 简化判断
		ID int `uri:"id" binding:"required,gt=1"`
	}{}
	// 2. 绑定
	result.Result(c.ShouldBindUri(id)).Unwrap() // 如果出错会发生panic然后被中间件捕捉
	// 3. 取值
	R(c)(
		"get_user_detail",
		"10001",
		getter.UserGetter.GetUserByID(id.ID).Unwrap(),
	)(OK)
}
