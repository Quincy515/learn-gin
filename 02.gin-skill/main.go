package main

import (
	"ginskill/src/common"
	"ginskill/src/dbs"
	"ginskill/src/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	dbs.InitDB()

	r := gin.New()
	r.Use(common.ErrorHandler())
	r.GET("/users", handlers.UserList)
	r.Run(":8080")
}
