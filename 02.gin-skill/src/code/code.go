package code

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()
	r.GET("/users", Handler()(GetUserList))
	r.GET("/users/:id", Handler()(GetUserDetail))
	r.Run(":8080")
}

type MyHandler func(c *gin.Context) (string, int, interface{})

func Handler() func(h MyHandler) gin.HandlerFunc {
	return func(h MyHandler) gin.HandlerFunc {
		return func(context *gin.Context) {
			msg, code, result := h(context)
			if code > 200 { // 代表可以做业务判断
				context.JSON(200, gin.H{"message": msg, "code": code, "result": result})
			} else {
				context.JSON(400, gin.H{"message": msg, "code": code, "result": result})
			}
		}
	}
}
func GetUserList(c *gin.Context) (string, int, interface{}) {
	// 各种业务代码
	return "userlist", 1001, "user_list"
}

func GetUserDetail(c *gin.Context) (string, int, interface{}) {
	// 各种业务代码
	return "userdetail", 1001, "user_detail:" + c.Param("id")
}
