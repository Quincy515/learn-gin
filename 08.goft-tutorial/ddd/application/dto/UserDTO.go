package dto

import "time"

// 输入对象
type (
	SimpleUserReq struct {
		Id int `uri:"id" binding:"required,min=100"`
	}
)

// 以下是输出对象
type (
	SimpleUserInfo struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		City string `json:"city"`
	}

	UserLog struct {
		Id   int       `json:"id"`
		Log  string    `json:"log"`
		Date time.Time `json:"date"`
	}

	UserInfo struct {
		Id    int        `json:"id"`
		Name  string     `json:"name"`
		City  string     `json:"city"`
		Phone string     `json:"phone"`
		Logs  []*UserLog `json:"logs"`
	}
)
