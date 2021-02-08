package main

import (
	"goft-tutorial/ddd/interfaces/configs"
	"goft-tutorial/ddd/interfaces/controllers"
	"goft-tutorial/pkg/goft"
)

func main() {
	goft.Ignite().
		Config(configs.NewUserServiceConfig(), configs.NewDBConfig(), configs.NewRepoConfig()).
		Mount("v1", controllers.NewUserController()).
		Launch()
}
