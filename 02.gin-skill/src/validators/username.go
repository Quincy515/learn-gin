package validators

import (
	"github.com/go-playground/validator/v10"
	"log"
)

func init() {
	if err := myvalid.RegisterValidation("UserName", VUserName); err != nil {
		log.Fatal("validator UserName error")
	}
}

var VUserName validator.Func = func(fl validator.FieldLevel) bool {
	uname, ok := fl.Field().Interface().(string) // 断言
	if ok && len(uname) >= 4 {
		return true
	}
	return false // true 表示验证通过，false 表示验证不通过
}
