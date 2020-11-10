package goft

import "github.com/gin-gonic/gin"

// Goft 嵌套 *gin.Engine
type Goft struct {
	*gin.Engine
}

// Ignite Goft 的构造函数，发射、燃烧，富含激情的意思
func Ignite() *Goft {
	return &Goft{Engine: gin.New()}
}

// Launch 最终启动函数，相当于 r.Run()
func (this *Goft) Launch() {
	this.Run(":8080")
}

// Mount 挂载控制器，定义接口，控制器继承接口就可以传进来
func (this *Goft) Mount(classes ...IClass) *Goft {
	for _, class := range classes {
		class.Build(this)
	}
	return this
}
