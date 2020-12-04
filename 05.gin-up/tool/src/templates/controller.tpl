package classes
import (
    	"github.com/gin-gonic/gin"
    	"gorm.io/gorm"
    	"log"
    	"gin-up/src/goft"
        "gin-up/src/models"
    )
{{$ClassName:=(printf "%s%s" .ControllerName "Class") | Ucfirst}}
type {{$ClassName}} struct { //控制器名称
		  *goft.GormAdapter  //注入Gorm 默认
}
func New{{$ClassName}}() *{{$ClassName}} {
	return &{{$ClassName}}{}
}
func(this *{{$ClassName}}) {{.ControllerName}}Detail(ctx *gin.Context) goft.Model{
	//obj:=models.New{{.ControllerName}}Model()
	//goft.Error(ctx.ShouldBindUri(obj))
	goft.Error(this.Table("your tablename").Where("id=?",11).Find(obj).Error)
	return obj
}
func(this *{{$ClassName}})  Build(goft *goft.Goft){
	//goft.Handle("GET","/your path/:id",this.{{.ControllerName}}Detail)
}
func(this *{{$ClassName}})  Name() string {
	 return "{{$ClassName}}"
}
