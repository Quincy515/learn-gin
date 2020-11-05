package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	. "topic/src"
)

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:root1234@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	tc := &TopicClass{}
	db.First(&tc, 2)
	fmt.Println(tc)

	var tcs []TopicClass
	db.Find(&tcs)
	fmt.Println(tcs)

	db.Where("class_name=?", "技术类").Find(&tcs)
	fmt.Println(tcs)
	db.Find(&tcs, "class_name=?", "新闻类")
	fmt.Println(tcs)
	db.Where(&TopicClass{ClassName: "技术类"}).Find(&tcs)
	fmt.Println(tcs)

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
