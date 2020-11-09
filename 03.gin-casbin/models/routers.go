package models

import "fmt"

type Routers struct {
	RouterId     int    `gorm:"column:r_id;primaryKey"`
	RouterName   string `gorm:"column:r_name"`
	RouterUri    string `gorm:"column:r_uri"`
	RouterMethod string `gorm:"column:r_method"`
	RoleName     string
	Domain       string `gorm:"column:tenant_name"`
}

func (this *Routers) TableName() string {
	return "roles"
}

func (this *Routers) String() string {
	return fmt.Sprintf("%s-%s", this.RouterName, this.RoleName)
}
