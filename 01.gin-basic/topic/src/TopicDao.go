package src

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// MustLogin 必须登录
func MustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, status := c.GetQuery("token"); !status {
			c.String(http.StatusUnauthorized, "缺少 token 参数")
			c.Abort()
		} else {
			c.Next()
		}
	}
}
func GetTopicDetail(c *gin.Context) {
	c.String(200, "获取topicid=%s的帖子", c.Param("topic_id"))
}

func NewTopic(c *gin.Context) {
	// 判断登录
	c.String(200, "新增帖子")
}

func DeleteTopic(c *gin.Context) {
	// 判断登录
	c.String(200, "删除帖子")
}
func GetTopicList() {

}
