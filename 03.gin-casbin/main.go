package main

import (
	"gin-casbin/lib"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(lib.Middlewares()...)

	r.GET("/:domain/depts", func(c *gin.Context) {
		c.JSON(200, gin.H{"result": "部门列表--" + c.Param("domain")})
	})
	r.POST("/:domain/depts", func(c *gin.Context) {
		c.JSON(200, gin.H{"reult": "批量修改部门列表" + c.Param("domain")})
	})
	r.Run(":8080")
}
