package LogModel

import "time"

// LogImpl 日志实体
type LogImpl struct {
	LogID   int       `gorm:"column:log_id,type:int,primaryKey,autoIncrement" json:"id"`
	LogName string    `json:"log_name"`
	LogDate time.Time `json:"log_date"`
}

func NewLogImpl(logName string, logDate time.Time) *LogImpl {
	return &LogImpl{LogName: logName, LogDate: logDate}
}

func (*LogImpl) TableName() string {
	return "logs"
}
