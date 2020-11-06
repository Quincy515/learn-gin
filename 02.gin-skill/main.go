package main

import (
	"ginskill/src/models/UserModel"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		user := UserModel.New()
		c.JSON(200, user)
	})
	r.Run(":8080")
}
