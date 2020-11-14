package main

import (
	. "gin-up/src/classes"
	. "gin-up/src/goft"
	. "gin-up/src/middlewares"
)

func main() {
	//GenTplFunc("src/funcs") // 在该参数目录下自动生成 funcmap.go 文件
	//return
	Ignite().
		Beans(NewGormAdapter(),NewXormAdapter()). // 设定数据库 orm 的 Bean，简单的依赖注入
		Attach(NewUserMid()). // 带声明周期的中间件
		Mount("v1", NewIndexClass(), // 控制器，挂载到 v1
			NewUserClass(), NewArticleClass()).
		Mount("v2", NewIndexClass()). // 控制器，挂载到 v2
		Task("0/3 * * * * *", Expr(".ArticleClass.Test")). // 每隔3秒，执行事件
		Launch()
}
