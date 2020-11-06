package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	. "topic/src"
)

func main2() {
	count := 0
	go func() { // 死循环程序
		for {
			fmt.Println("执行", count)
			count++
			time.Sleep(time.Second * 1)
		}
	}()

	c := make(chan os.Signal) // 1. 创建信号 chan
	go func() {
		// 2. 创建一个 超时 context，到期后会执行 Done
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		// 中间是业务
		select {
		case <-ctx.Done():
			c <- os.Interrupt // 3. 重点，超时时间到了会发送 SIGINT 信号
		}
	}()
	signal.Notify(c) // 4. 监听所有信号
	s := <-c         // 赋值给变量 s
	fmt.Println(s)
}
func main() {
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

	//router.Run() // 8080
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go (func() { // 启动 web 服务
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("服务器启动失败:", err)
		}
	})()
	go (func() { // 通过协程方式启动数据库
		InitDB()
	})()
	ServerNotify() // 信号监听
	//这里还可以做一些 释放连接或善后工作，暂时略
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil { // 强制关闭
		log.Fatalln("服务器关闭")
	}
	log.Println("服务器优雅退出")
}
