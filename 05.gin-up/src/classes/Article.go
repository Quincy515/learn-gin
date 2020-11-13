package classes

import (
	"gin-up/src/goft"
	"gin-up/src/models"
	"github.com/gin-gonic/gin"
)

type ArticleClass struct {
	*goft.GormAdapter
}

func NewArticleClass() *ArticleClass {
	return &ArticleClass{}
}

func (this *ArticleClass) ArticleDetail(ctx *gin.Context) goft.Model {
	news := models.NewArticleModel()
	goft.Error(ctx.ShouldBindUri(news))
	res := this.Table(news.TableName()).Where("id=?", news.NewsId).Find(news)
	if res.Error != nil || res.RowsAffected == 0 {
		goft.Error(errors.New("没有找到记录"))
	}
	return news
}

func (this *ArticleClass) Build(goft *goft.Goft) {
	goft.Handle("GET", "/article/:id", this.ArticleDetail)
}
