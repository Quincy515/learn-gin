package main

import (
	"ginskill/src/common"
	"ginskill/src/models/UserModel"
	"ginskill/src/result"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(common.ErrorHandler())
	r.POST("/", func(c *gin.Context) {
		user := UserModel.New()
		result.Result(c.ShouldBind(user)).Unwrap()
		c.JSON(200, user)
	})
	r.Run(":8080")
}
