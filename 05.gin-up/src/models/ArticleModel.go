package models

import (
	"fmt"
	"time"
)

type ArticleModel struct {
	NewsId      int    `json:"id" gorm:"column:id" uri:"id" binding:"required,gt=0"`
	Newstitle   string `json:"title"`
	Newscontent string `json:"content"`
	Views       int    `json:"views"`
	Addtime     time.Time
}

func NewArticleModel() *ArticleModel {
	return &ArticleModel{}
}

func (ArticleModel) TableName() string {
	return "mynews"
}

func (this *ArticleModel) String() string {
	return fmt.Sprintf("id: %d, title: %s", this.NewsId, this.Newstitle)
}
