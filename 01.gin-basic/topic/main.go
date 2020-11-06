package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	. "topic/src"
)

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	//dsn := "root:root1234@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	//db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//topic := Topic{
	//	TopicTitle:      "TopicTitle",
	//	TopicShortTitle: "TopicShortTitle",
	//	UserIP:          "127.0.0.1",
	//	TopicScore:      0,
	//	TopicUrl:        "testurl",
	//	TopicDate:       time.Now()}
	//result := db.Create(&topic)
	//fmt.Println(topic.TopicID)
	//fmt.Println(result.Error)
	//fmt.Println(result.RowsAffected)

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("topicurl", TopicUrl)
		v.RegisterValidation("topics", TopicsValidate) // 自定义验证批量post帖子的长度
	}

	v1 := router.Group("/v1/topics") // 单条帖子处理
	{
		v1.GET("", GetTopicList)

		v1.GET("/:topic_id", GetTopicDetail)

		v1.Use(MustLogin())
		{
			v1.POST("", NewTopic)
			v1.DELETE("/:topic_id", DeleteTopic)
		}
	}

	v2 := router.Group("/v1/mtopics") // 多条帖子处理
	{
		v2.Use(MustLogin())
		{
			v2.POST("", NewTopics)
		}
	}

	router.Run() // 8080
}
