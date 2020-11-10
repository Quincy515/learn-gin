package classes

import "github.com/gin-gonic/gin"

// UserClass *gin.Engine 的嵌套
type UserClass struct {
	*gin.Engine
}

// NewUserClass UserClass generate constructor
func NewUserClass(engine *gin.Engine) *UserClass {
	return &UserClass{Engine: engine}
}

// UserList 控制器方法
func (this *UserClass) UserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": "success",
		})
	}
}

func (this *UserClass) Build() {
	this.Handle("GET", "/user", this.UserList())
}
