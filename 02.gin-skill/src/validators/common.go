package validators

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
)

var myvalid *validator.Validate
var validatorError map[string]string

func init() {
	// 包初始化时 make，所以在运行过程中不会增加或减少，所以是线程安全的。
	validatorError = make(map[string]string)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		myvalid = v
	} else {
		log.Fatal("error validator")
	}
}

func registerValidation(tag string, fn validator.Func) {
	err := myvalid.RegisterValidation(tag, fn)
	if err != nil {
		log.Fatal(fmt.Sprintf("validator %s error", tag))
	}
}

// CheckErrors 断言是否是验证错误
func CheckErrors(errors error) {
	if errs, ok := errors.(validator.ValidationErrors); ok {
		for _, err := range errs { // 如果是验证错误，看有没有自定义的错误信息
			if v, exists := validatorError[err.Tag()]; exists {
				panic(v)
			}
		}
	}
}

// tagErrMsg 默认的 tag 出错信息
func tagErrMsg() {
	validatorError["min"] = "位数太少"
}
