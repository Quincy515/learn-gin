package classes
import (
    	"github.com/gin-gonic/gin"
    	"gorm.io/gorm"
    	"log"
    	"gin-up/src/goft"
        "gin-up/src/models"
    )

type UserClass struct { //控制器名称
		  *goft.GormAdapter  //注入Gorm 默认
}
func NewUserClass() *UserClass {
	return &UserClass{}
}
func(this *UserClass) userDetail(ctx *gin.Context) goft.Model{
	//obj:=models.NewuserModel()
	//goft.Error(ctx.ShouldBindUri(obj))
	goft.Error(this.Table("your tablename").Where("id=?",11).Find(obj).Error)
	return obj
}
func(this *UserClass)  Build(goft *goft.Goft){
	//goft.Handle("GET","/your path/:id",this.userDetail)
}
func(this *UserClass)  Name() string {
	 return "UserClass"
}
