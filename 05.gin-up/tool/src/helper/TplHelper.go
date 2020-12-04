package Helper

import (
	"log"
	"os"
	"text/template"
)

// 通过模板引擎 生成文件
// tplContent 是模板内容
func GenFile(tplContent string, dest string, data map[string]interface{}) {
	tpl := template.New("goft-tpl").Funcs(NewTplFunction())
	tmpl, err := tpl.Parse(tplContent)
	if err != nil {
		log.Fatal(" tpl parse-error:", err)
	}
	file, err := os.OpenFile(GetWorkDir()+"/"+dest,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("generate target error:", err)
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatal("generate error:", err)
	}
}
