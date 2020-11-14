package goft

import "reflect"

type Bean interface {
	Name() string // 所有控制器加入到 bean 中，都需要增加 Name() 函数
}

type BeanFactory struct {
	beans []Bean
}

func NewBeanFactory() *BeanFactory {
	bf := &BeanFactory{beans: make([]Bean, 0)}
	bf.beans = append(bf.beans, bf)
	return bf
}

// setBean 实现简单的依赖注入，往内存中塞入 bean
func (this *BeanFactory) setBean(beans ...Bean) {
	this.beans = append(this.beans, beans...)
}

// GetBean 外部使用获取注入的属性
func (this *BeanFactory) GetBean(bean interface{}) interface{} {
	return this.getBean(reflect.TypeOf(bean))
}

// getBean 获取内存中预先设置好的 bean 对象
func (this *BeanFactory) getBean(t reflect.Type) interface{} {
	for _, p := range this.beans {
		if t == reflect.TypeOf(p) {
			return p
		}
	}
	return nil
}

// Inject 给外部用的 （后面还要改,这个方法不处理注解)
func (this *BeanFactory) Inject(object interface{}) {
	vObject := reflect.ValueOf(object)
	if vObject.Kind() == reflect.Ptr { //由于不是控制器 ，所以传过来的值 不一定是指针。因此要做判断
		vObject = vObject.Elem()
	}
	for i := 0; i < vObject.NumField(); i++ {
		f := vObject.Field(i)
		if f.Kind() != reflect.Ptr || !f.IsNil() {
			continue
		}
		if p := this.getBean(f.Type()); p != nil && f.CanInterface() {
			f.Set(reflect.New(f.Type().Elem()))
			f.Elem().Set(reflect.ValueOf(p).Elem())
		}
	}
}

// inject 把 bean 注入到控制器中 (内部方法,用户控制器注入。并同时处理注解)
func (this *BeanFactory) inject(class IClass) {
	vClass := reflect.ValueOf(class).Elem()
	vClassT := reflect.TypeOf(class).Elem()
	for i := 0; i < vClass.NumField(); i++ {
		f := vClass.Field(i)
		if f.Kind() != reflect.Ptr || !f.IsNil() {
			continue
		}
		if IsAnnotation(f.Type()) { // 如果是注解，单独处理
			f.Set(reflect.New(f.Type().Elem()))
			f.Interface().(Annotation).SetTag(vClassT.Field(i).Tag)
			this.Inject(f.Interface())
			continue
		}
		if p := this.getBean(f.Type()); p != nil {
			f.Set(reflect.New(f.Type().Elem()))
			f.Elem().Set(reflect.ValueOf(p).Elem())
		}
	}
}

func (this *BeanFactory) Name() string {
	return "BeanFactory"
}
