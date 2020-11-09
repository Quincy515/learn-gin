package mappers

import (
	"ginskill/src/dbs"
	"gorm.io/gorm"
)

type SqlMapper struct {
	Sql  string
	Args []interface{}
}

func NewSqlMapper(sql string, args []interface{}) *SqlMapper {
	return &SqlMapper{Sql: sql, Args: args}
}

// Mapper 转换返回值生成 SqlMapper
func Mapper(sql string, args []interface{}, err error) *SqlMapper {
	if err != nil {
		panic(err.Error())
	}
	return NewSqlMapper(sql, args)
}

// Query 对 SqlMapper 进行查询的封装
func (this *SqlMapper) Query() *gorm.DB {
	return dbs.Orm.Raw(this.Sql, this.Args)
}

// Exec 对 SqlMapper 进行执行 update/delete/inset 的封装
func (this *SqlMapper) Exec() *gorm.DB {
	return dbs.Orm.Exec(this.Sql, this.Args...)
}
