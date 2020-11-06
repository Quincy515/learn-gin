package validators

import (
	"github.com/go-playground/validator/v10"
)

// init 注册自定义验证
func init() {
	// tag 规则强制转换为自定义的规则名 UserName
	registerValidation("UserName", UserName("required,min=4").toFunc())
}

// UserName 就是规则
type UserName string // 用户名是 string 类型

// 针对用户名的验证规则
func (this UserName) toFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		v, ok := fl.Field().Interface().(string) // 断言
		if ok {
			return this.validate(v)
		}
		return false
	}
}

func (this UserName) validate(v string) bool {
	// 本身的 tag 验证
	if err := myvalid.Var(v, string(this)); err != nil { // 单字段验证
		return false
	}
	// 其他自定义验证
	if len(v) > 8 {
		return false
	}
	return true
}

//var VUserName validator.Func = func(fl validator.FieldLevel) bool {
//	uname, ok := fl.Field().Interface().(string) // 断言
//	if ok && len(uname) >= 4 {
//		return true
//	}
//	return false // true 表示验证通过，false 表示验证不通过
//}
