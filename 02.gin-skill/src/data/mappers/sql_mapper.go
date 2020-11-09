package mappers

import (
	"ginskill/src/dbs"
	"gorm.io/gorm"
)

type SqlMapper struct {
	Sql  string        // sql 语句
	Args []interface{} // 参数集合
	db   *gorm.DB
}

// setDB 私有方法
func (this *SqlMapper) setDB(db *gorm.DB) {
	this.db = db
}

func NewSqlMapper(sql string, args []interface{}) *SqlMapper {
	return &SqlMapper{Sql: sql, Args: args} // 单表执行
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
	if this.db != nil { // 不是单表执行
		return this.db.Raw(this.Sql, this.Args...)
	}
	return dbs.Orm.Raw(this.Sql, this.Args...) // 单表执行
}

// Exec 对 SqlMapper 进行执行 update/delete/inset 的封装
func (this *SqlMapper) Exec() *gorm.DB {
	if this.db != nil { // 不是单表执行
		return this.db.Exec(this.Sql, this.Args...)
	}
	return dbs.Orm.Exec(this.Sql, this.Args...) // 单表执行
}

// SqlMappers 定义多表的事物操作
type SqlMappers []*SqlMapper

func Mappers(sqlMappers ...*SqlMapper) (list SqlMappers) {
	list = sqlMappers
	return
}

// 执行 Mappers 方法，就把所有 sql 执行设置为同一个 db
func (this SqlMappers) apply(tx *gorm.DB) {
	for _, sql := range this {
		sql.setDB(tx) // 多表执行
	}
}

// Exec 多表执行，传入函数 f，表示执行事物
func (this SqlMappers) Exec(f func() error) error {
	return dbs.Orm.Transaction(func(tx *gorm.DB) error {
		this.apply(tx) // tx 就是统一使用的 db 对象
		return f()     // gorm Transaction 机制返回的是 error 就自动回滚，不是 error 就执行 commit 操作
	})
}
