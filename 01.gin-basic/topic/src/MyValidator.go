package src

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func TopicsValidate(fl validator.FieldLevel) bool {
	topics, ok := fl.Top().Interface().(*Topics)
	if ok && topics.TopicListSize == len(topics.TopicList) {
		return true
	}
	return false
}

// TopicUrl 自定义字段级别校验方法
func TopicUrl(fl validator.FieldLevel) bool {
	// 判断当前传入的 struct 是否是 Topic model
	_, ok1 := fl.Top().Interface().(*Topic)
	_, ok2 := fl.Top().Interface().(*Topics)
	if ok1 || ok2 {
		getValue := fl.Field().String()
		if ret, _ := regexp.MatchString("^\\w{4,10}$", getValue); ret {
			return true
		}
	}
	return false
}
