package models

import "time"

const (
	UserLog_Create = 5
	UserLog_Update = 6
)

// UserLogModel 用户日志实体
type UserLogModel struct {
	*Model
	Id         int       `gorm:"column:id;primary_key;auto_increment" json:"id"`
	UserName   string    `gorm:"column:user_name" json:"user_name"`
	LogType    uint8     `gorm:"column:log_type" json:"log_type"`
	LogComment string    `gorm:"column:log_comment" json:"log_comment"`
	Updatetime time.Time `gorm:"column:update_time" json:"login_time"`
}

func NewUserLogModel(userName string, attrs ...UserLogAttrFunc) *UserLogModel {
	logModel := &UserLogModel{UserName: userName}
	UserLogAttrFuncs(attrs).apply(logModel)
	logModel.Model = &Model{}
	logModel.SetId(logModel.Id)
	logModel.SetName("用户日志实体")
	logModel.Updatetime = time.Now()
	return logModel
}
