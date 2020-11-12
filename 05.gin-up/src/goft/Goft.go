package goft

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

// Goft
type Goft struct {
	*gin.Engine                  // 把 *gin.Engine 放入主类里
	g           *gin.RouterGroup // 保存 group 对象
	props       []interface{}
	//dba         interface{}      // 保存和执行 *gorm.DB 对象
}

// Ignite Goft 的构造函数，发射、燃烧，富含激情的意思
func Ignite() *Goft {
	g := &Goft{Engine: gin.New(), props: make([]interface{}, 0)}
	g.Use(ErrorHandler()) // 必须强制加载异常处理中间件
	return g
}

// Launch 最终启动函数，相当于 r.Run()
func (this *Goft) Launch() {
	config := InitConfig()
	this.Run(fmt.Sprintf(":%d", config.Server.Port))
}

// Handle 重载 gin.Handle 函数
func (this *Goft) Handle(httpMethod, relativePath string, handler interface{}) *Goft {
	if h := Convert(handler); h != nil {
		this.g.Handle(httpMethod, relativePath, h)
	}
	return this
}

// Attach 实现中间件的加入
func (this *Goft) Attach(f Fairing) *Goft {
	this.Use(func(c *gin.Context) {
		err := f.OnRequest(c)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		} else {
			c.Next() // 继续往下执行
		}
	})
	return this
}

// Beans 实现简单的依赖注入
func (this *Goft) Beans(beans ...interface{}) *Goft {
	this.props = append(this.props, beans...)
	return this
}

// Mount 挂载控制器，定义接口，控制器继承接口就可以传进来
func (this *Goft) Mount(group string, classes ...IClass) *Goft {
	this.g = this.Group(group)
	for _, class := range classes {
		class.Build(this)
		this.setProp(class)
	}
	return this
}

// getProp 获取属性
func (this *Goft) getProp(t reflect.Type) interface{} {
	for _, p := range this.props {
		if t == reflect.TypeOf(p) {
			return p
		}
	}
	return nil
}

// setProp 赋值
func (this *Goft) setProp(class IClass) {
	// reflect.ValueOf(class) 是指针，reflect.ValueOf(class) 是指针指向的对象
	vClass := reflect.ValueOf(class).Elem()  // 反射
	vClassT := reflect.TypeOf(class).Elem()
	for i := 0; i < vClass.NumField(); i++ { // 遍历 vClass 的属性
		f := vClass.Field(i)                       // 判断属性是否已经初始化
		if !f.IsNil() || f.Kind() != reflect.Ptr { // 如果控制器已经初始化或者不是指针
			continue // 就跳过
		}
		if p := this.getProp(f.Type()); p != nil {
			// vClass.Field(0)是强制使用第一个属性的指针，使用 Set() 进行赋值完成初始化
			// vClass.Field(0).Type() --> 指针 *GormAdapter
			// vClass.Field(0).Type().Elem() -->指针指向的对象 GormAdapter
			// reflect.New(vClass.Field(0).Type().Elem()) --> new 指针 *GormAdapter
			// vClass.Field(0).Set(reflect.New(vClass.Field(0).Type().Elem()))
			// Elem() 是指针指向的对象 Set() 是进行赋值
			// vClass.Field(0).Elem().Set(reflect.ValueOf(this.dba).Elem())
			f.Set(reflect.New(f.Type().Elem()))     // 初始化
			f.Elem().Set(reflect.ValueOf(p).Elem()) // 赋值

			if IsAnnotation(f.Type()) { // 判断是否是注解
				p.(Annotation).SetTag(vClassT.Field(i).Tag)
			}
		}
	}
}
