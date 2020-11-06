package main

import (
	"ginskill/src/common"
	"ginskill/src/dbs"
	"ginskill/src/handlers"
	_ "ginskill/src/validators" // 不调用，仅引用，执行 init 函数即可
	"github.com/gin-gonic/gin"
)

func main() {
	dbs.InitDB()

	r := gin.New()
	r.Use(common.ErrorHandler())
	r.GET("/users", handlers.UserList)
	r.GET("/users/:id", handlers.UserDetail)
	r.POST("/users", handlers.UserSave)
	r.Run(":8080")
}
