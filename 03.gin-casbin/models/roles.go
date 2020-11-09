package models

import "fmt"

type Role struct {
	RoleId      int    `gorm:"column:role_id"`
	RoleName    string `gorm:"column:role_name"`
	RolePid     int    `gorm:"column:role_pid"`
	RoleComment string `gorm:"column:role_comment"`
	TenantId    string `gorm:"column:tenant_id"`
	TenantName  string `gorm:"column:tenant_name"`
}

func (this *Role) TableName() string {
	return "roles"
}

func (this *Role) String() string {
	return fmt.Sprintf("ID:%d 角色名:%s - 租户名:%s", this.RoleId, this.RoleName, this.TenantName)
}
