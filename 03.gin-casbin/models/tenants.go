package models

import "fmt"

// 租户模型
type Tenant struct {
	TenantId   int    `gorm:"column:tenant_id;primaryKey"`
	TenantName string `gorm:"column:tenant_name"`
}

func (this *Tenant) String() string {
	return fmt.Sprintf("%d:%s", this.TenantId, this.TenantName)
}
