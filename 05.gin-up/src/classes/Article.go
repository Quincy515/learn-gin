package classes

import (
	"errors"
	"gin-up/src/goft"
	"gin-up/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"time"
)

type ArticleClass struct {
	*goft.GormAdapter
}

func NewArticleClass() *ArticleClass {
	return &ArticleClass{}
}

// ArticleDetail 路由到这个方法里面
func (this *ArticleClass) ArticleDetail(ctx *gin.Context) goft.Model {
	news := models.NewArticleModel()
	goft.Error(ctx.ShouldBindUri(news))
	res := this.Table(news.TableName()).Where("id=?", news.NewsId).Find(news)
	if res.Error != nil || res.RowsAffected == 0 {
		goft.Error(errors.New("没有找到记录"))
	}

	//goft.Task(this.UpdateView, this.UpdateViewDone, news.NewsId) // 代表执行一个协程任务
	goft.Task(this.UpdateView, func() { // 通过匿名函数带参数
		this.UpdateViewDoneWithParam(news.NewsId)
	}, news.NewsId) // 代表执行一个协程任务

	return news
}

// UpdateView 增加点击量
func (this *ArticleClass) UpdateView(params ...interface{}) {
	time.Sleep(time.Second * 3)
	this.Table("mynews").Where("id=?", params[0]).
		Update("views", gorm.Expr("views+1"))
}

// UpdateViewDone 测试回调函数
func (this *ArticleClass) UpdateViewDone() {
	log.Println("点击量增加结束")
}

// UpdateViewDoneWithParam 测试带参数的回调函数
func (this *ArticleClass) UpdateViewDoneWithParam(id int) {
	log.Println("点击量增加结束, id 是: ", id)
}

// Test 控制器里的定时任务方法
func (this *ArticleClass) Test() interface{} {
	log.Println("测试定时任务")
	return nil
}

func (this *ArticleClass) Build(goft *goft.Goft) {
	goft.Handle("GET", "/article/:id", this.ArticleDetail)
}

// Name 控制器加入到 bean 中
func (this *ArticleClass) Name() string {
	return "ArticleClass"
}
