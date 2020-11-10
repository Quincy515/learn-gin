package main

import (
	. "gin-up/src/classes"
	"gin-up/src/goft"
)

func main() {
	goft.Ignite().Mount(NewIndexClass(), NewUserClass()).Launch()
}
