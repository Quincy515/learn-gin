package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:root1234@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	rows, _ := db.Raw("SELECT topic_id, topic_title FROM topics").Rows()
	for rows.Next() {
		var t_id int
		var t_title string
		rows.Scan(&t_id, &t_title)
		fmt.Println(t_id, t_title)
	}
	//router := gin.Default()
	//
	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//	v.RegisterValidation("topicurl", TopicUrl)
	//	v.RegisterValidation("topics", TopicsValidate) // 自定义验证批量post帖子的长度
	//}
	//
	//v1 := router.Group("/v1/topics") // 单条帖子处理
	//{
	//	v1.GET("", GetTopicList)
	//
	//	v1.GET("/:topic_id", GetTopicDetail)
	//
	//	v1.Use(MustLogin())
	//	{
	//		v1.POST("", NewTopic)
	//		v1.DELETE("/:topic_id", DeleteTopic)
	//	}
	//}
	//
	//v2 := router.Group("/v1/mtopics") // 多条帖子处理
	//{
	//	v2.Use(MustLogin())
	//	{
	//		v2.POST("", NewTopics)
	//	}
	//}
	//
	//router.Run() // 8080
}
