package src

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	DBHelper *gorm.DB
	err      error
)

func InitDB() {
	dsn := "root:root1234@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	DBHelper, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		//fmt.Println(err)
		//log.Fatal("数据库初始化错误：", err)
		ShutDownServer(err)
		return
	}
	sqlDB, _ := DBHelper.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}
