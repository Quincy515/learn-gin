package main

import (
	. "gin-up/src/classes"
	"gin-up/src/goft"
	. "gin-up/src/middlewares"
)

func main() {
	goft.Ignite().
		Attach(NewUserMid()). // 带声明周期的中间件
		Mount("v1", NewIndexClass(), // 控制器，挂载到 v1
			NewUserClass()).
		Mount("v2", NewIndexClass()). // 控制器，挂载到 v2
		Launch()
}
