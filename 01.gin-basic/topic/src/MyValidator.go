package src

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// TopicUrl 自定义字段级别校验方法
func TopicUrl(fl validator.FieldLevel) bool {
	// 判断当前传入的 struct 是否是 Topic model
	if _, ok := fl.Top().Interface().(*Topic); ok {
		getValue := fl.Field().String()
		if ret, _ := regexp.MatchString("^\\w{4,10}$", getValue); ret {
			return true
		}
	}
	return false
}
