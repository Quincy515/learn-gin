package models

import "time"

// UserLogModel 用户日志实体
type UserLogModel struct {
	*Model
	Id         int       `gorm:"column:id;primary_key;auto_increment" json:"id"`
	UserName   string    `gorm:"column:user_name" json:"user_name"`
	LogType    uint8     `gorm:"column:log_type" json:"log_type"`
	LogComment uint8     `gorm:"column:log_comment" json:"log_comment"`
	Updatetime time.Time `gorm:"column:update_time" json:"login_time"`
}

func NewUserLogModel(userName string, logType uint8, logComment uint8) *UserLogModel {
	logModel := &UserLogModel{UserName: userName, LogType: logType, LogComment: logComment}
	logModel.Model = &Model{}
	logModel.SetId(logModel.Id)
	logModel.SetName("用户日志实体")
	return logModel
}
