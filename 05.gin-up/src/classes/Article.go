package classes

import (
	"errors"
	"gin-up/src/goft"
	"gin-up/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	goft.Task(this.UpdateView, news.NewsId) // 代表执行一个协程任务
	return news
}

// UpdateView 增加点击量
func (this *ArticleClass) UpdateView(params ...interface{}) {
	this.Table("mynews").Where("id=?", params[0]).
		Update("views", gorm.Expr("views+1"))
}

func (this *ArticleClass) Build(goft *goft.Goft) {
	goft.Handle("GET", "/article/:id", this.ArticleDetail)
}
