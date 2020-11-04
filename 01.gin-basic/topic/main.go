package main

import (
	"github.com/gin-gonic/gin"
	. "topic/src"
)

func main() {
	router := gin.Default()

	v1 := router.Group("/v1/topics")
	{
		v1.GET("", func(c *gin.Context) {
			if c.Query("username") == "" {
				c.String(200, "获取帖子列表")
			} else {
				c.String(200, "获取用户=%s的帖子列表", c.Query("username"))
			}
		})

		v1.GET("/:topic_id", GetTopicDetail)

		v1.Use(MustLogin())
		{
			v1.POST("", NewTopic)
			v1.DELETE("/:topic_id", DeleteTopic)
		}
	}

	router.Run() // 8080
}
