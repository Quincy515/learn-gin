package models

type UserModel struct {
	UserID       int `gorm:"column:user_id"`
	UserName     string
	UserPwd      string
	SourceID     string `gorm:"column:source_id"`
	SourceUserId string `gorm:"column:source_userid"`
}

func (*UserModel) TableName() string {
	return "users"
}
