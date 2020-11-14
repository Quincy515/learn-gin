package classes

import (
	"gin-up/src/goft"
	"github.com/gin-gonic/gin"
)

type IndexClass struct{}

func NewIndexClass() *IndexClass {
	return &IndexClass{}
}

// GetIndex 业务方法，函数名根据业务而起
func (i *IndexClass) GetIndex(ctx *gin.Context) goft.View {
	ctx.Set("name", "custer")
	return "index"
}

// Build 把业务的路由隐藏在 Build 函数
func (i *IndexClass) Build(goft *goft.Goft) {
	goft.Handle("GET", "/", i.GetIndex)
}

func (this *IndexClass) Name() string {
	return "IndexClass"
}
