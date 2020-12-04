package goft

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

type GormAdapter struct {
	*gorm.DB
}

func (this *GormAdapter) Name() string {
	return "GormAdapter"
}
func NewGormAdapter() *GormAdapter {
	dsn := "root:root1234@tcp(127.0.0.1:3306)/casbin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		//TablePrefix: "t_",   // table name prefix, table for `User` would be `t_users`
		SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
	}})
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB.SetMaxIdleConns(5)                   //最大空闲数
	mysqlDB.SetMaxOpenConns(10)                  //最大打开连接数
	mysqlDB.SetConnMaxLifetime(time.Second * 30) //空闲连接生命周期
	return &GormAdapter{DB: db}
}
