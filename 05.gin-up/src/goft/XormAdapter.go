package goft

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"xorm.io/xorm"
)

type XormAdapter struct {
	*xorm.Engine
}

func (this *XormAdapter) Name() string {
	return "XormAdapter"
}

func NewXormAdapter() *XormAdapter {
	dsn := "root:root1234@tcp(127.0.0.1:3306)/client?charset=utf8mb4&parseTime=True&loc=Local"
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := engine.Ping(); err == nil {
		log.Println("通过 xorm 连接数据库成功")
	}
	engine.DB().SetMaxIdleConns(5)  // SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	engine.DB().SetMaxOpenConns(10) // SetMaxOpenConns 设置打开数据库连接的最大数量。
	return &XormAdapter{Engine: engine}
}
