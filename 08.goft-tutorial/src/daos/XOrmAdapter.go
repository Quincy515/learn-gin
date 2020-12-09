package daos

import (
	"database/sql"
	"xorm.io/xorm"
)

type XOrmAdapter struct {
	*xorm.Engine
}

func (this *XOrmAdapter) DB() *sql.DB {
	return this.Engine.DB().DB
}
