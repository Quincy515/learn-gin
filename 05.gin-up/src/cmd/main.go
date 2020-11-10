package main

import (
	. "gin-up/src/classes"
	"gin-up/src/goft"
)

func main() {
	goft.Ignite().
		Mount("v1", NewIndexClass(),
			NewUserClass()).
		Mount("v2", NewIndexClass()).
		Launch()
}
