package main

import (
	"goft-tutorial/pkg/goft"
	"goft-tutorial/src/configure"
	"goft-tutorial/src/controllers"
	"goft-tutorial/src/middlewares"
)

func main() {
	goft.Ignite().
		Config(configure.NewDBConfig(), configure.NewServiceConfig()).
		Attach(middlewares.NewAddVersion()).
		Mount("v1", controllers.NewIndexController(),
			controllers.NewUserController()).
		Launch()
}
