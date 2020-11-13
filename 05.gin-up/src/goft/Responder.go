package goft

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

// ResponderList 切片是一堆 Responder 接口
var ResponderList []Responder

func init() {
	ResponderList = []Responder{
		new(StringResponder),
		new(ModelResponder),
		new(ModelsResponder),
		new(ViewResponder),
	} // 反射不能直接使用类型，提供反射需要的指针
}

// Responder 接口目的是 gin.Handle() 的第三个参数
type Responder interface {
	RespondTo() gin.HandlerFunc
}

// Convert 通过反射判断 interface{] 类型断言
func Convert(handler interface{}) gin.HandlerFunc {
	hRef := reflect.ValueOf(handler) // handler 变成 reflect 反射对象
	for _, r := range ResponderList {
		rRef := reflect.ValueOf(r).Elem() // new() 的指针类型必须要执行 Elem
		// 判断 hRef 的类型是否可以转换成 rRef 的类型
		if hRef.Type().ConvertibleTo(rRef.Type()) {
			rRef.Set(hRef) // 反射的方式设置值
			return rRef.Interface().(Responder).RespondTo()
		}
	}
	return nil
}

// StringResponder 把返回字符串的 gin.HandlerFunc 包装成一个类型
type StringResponder func(*gin.Context) string

// RespondTo 接口的实现
func (this StringResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(200, this(context))
	}
}

// ModelResponder 控制器返回实体类
type ModelResponder func(*gin.Context) Model

// RespondTo 接口的实现
func (this ModelResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, this(context))
	}
}

// ModelsResponder 控制器返回实体类切片
type ModelsResponder func(*gin.Context) Models

// RespondTo 接口的实现
func (this ModelsResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Content-type", "application/json")
		context.Writer.WriteString(string(this(context)))
	}
}

// View
type View string

// ViewResponder 返回 html 模板
type ViewResponder func(*gin.Context) View

// RespondTo 接口的实现
func (this ViewResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.HTML(200, string(this(context))+".html", nil)
	}
}
