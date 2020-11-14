package goft

import (
	"fmt"
	"gin-up/src/funcs"
	"github.com/gin-gonic/gin"
	"log"
)

// Goft
type Goft struct {
	*gin.Engine                  // 把 *gin.Engine 放入主类里
	g           *gin.RouterGroup // 保存 group 对象
	beanFactory *BeanFactory
	exprData    map[string]interface{} // 表达式的数据
}

// Ignite Goft 的构造函数，发射、燃烧，富含激情的意思
func Ignite() *Goft {
	g := &Goft{Engine: gin.New(), beanFactory: NewBeanFactory(), exprData: map[string]interface{}{}}
	g.Use(ErrorHandler()) // 必须强制加载异常处理中间件
	config := InitConfig()
	g.beanFactory.setBean(config) // 配置文件加载进 bean 中
	if config.Server.Html != "" {
		g.FuncMap = funcs.FuncMap
		g.LoadHTMLGlob(config.Server.Html)
	}
	return g
}

// Launch 最终启动函数，相当于 r.Run()
func (this *Goft) Launch() {
	var port int32 = 8080
	if config := this.beanFactory.GetBean(new(SysConfig)); config != nil {
		port = config.(*SysConfig).Server.Port
	}
	getCronTask().Start()
	this.Run(fmt.Sprintf(":%d", port))
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

// Beans 加入 bean 容器
func (this *Goft) Beans(beans ...Bean) *Goft {
	// 取出 bean 的名称，加入到 exprData 里面
	for _, bean := range beans {
		this.exprData[(bean.Name())] = bean
	}
	this.beanFactory.setBean(beans...)
	return this
}

// Mount 挂载控制器，定义接口，控制器继承接口就可以传进来
func (this *Goft) Mount(group string, classes ...IClass) *Goft {
	this.g = this.Group(group)
	for _, class := range classes {
		class.Build(this)
		this.beanFactory.inject(class)
		this.Beans(class) // 控制器也作为 bean 加入到 bean 容器
	}
	return this
}

// Task 增加定时任务 参数 cron 表示式
func (this *Goft) Task(cron string, expr interface{}) *Goft {
	var err error
	if f, ok := expr.(func()); ok { // 断言 expr 是一个 func
		_, err = getCronTask().AddFunc(cron, f)
	} else if exp, ok := expr.(Expr); ok { // 断言 expr 是一个 表达式类型
		_, err = getCronTask().AddFunc(cron, func() {
			_, expErr := ExecExpr(exp, this.exprData)
			if expErr != nil {
				log.Println(expErr)
			}
		})
	}

	if err != nil {
		log.Println(err)
	}
	return this
}
