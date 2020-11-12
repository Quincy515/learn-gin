package goft

import (
	"fmt"
	"reflect"
	"strings"
)

// Annotation 注解接口
type Annotation interface {
	SetTag(tag reflect.StructTag) // 通过 Tag 完成更多复杂功能
	String() string
}

// AnnotationList 注解列表是注解接口的切片
var AnnotationList []Annotation

//IsAnnotation 判断当前注入对象是否是注解，运行在系统启动之前运行，不用考虑性能
func IsAnnotation(t reflect.Type) bool {
	for _, item := range AnnotationList {
		if reflect.TypeOf(item) == t {
			return true
		}
	}
	return false
}

// init 包构造函数
func init() {
	AnnotationList = make([]Annotation, 0)
	AnnotationList = append(AnnotationList, new(Value))
}

// Value 注解
type Value struct {
	tag         reflect.StructTag
	BeanFactory *BeanFactory
}

func (this *Value) SetTag(tag reflect.StructTag) {
	this.tag = tag
}

func (this *Value) String() string {
	get_prefix := this.tag.Get("prefix")
	if get_prefix == "" {
		return ""
	}
	prefix := strings.Split(get_prefix, ".")
	if config := this.BeanFactory.GetBean(new(SysConfig)); config != nil {
		get_value := GetConfigValue(config.(*SysConfig).Config, prefix, 0)
		if get_value != nil {
			return fmt.Sprintf("%v", get_value)
		} else {
			return ""
		}
	} else {
		return ""
	}
}
