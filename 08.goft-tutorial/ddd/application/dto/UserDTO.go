package dto

import "time"

type SimpleUserInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type UserLog struct {
	Id   int       `json:"id"`
	Log  string    `json:"log"`
	Date time.Time `json:"date"`
}

type UserInfo struct {
	Id    int        `json:"id"`
	Name  string     `json:"name"`
	City  string     `json:"city"`
	Phone string     `json:"phone"`
	Logs  []*UserLog `json:"logs"`
}
