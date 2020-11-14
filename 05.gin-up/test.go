package main

import (
	"fmt"
	"gin-up/src/goft"
)

func main() {
	// 使用 go 原生的模板功能，
	//tpl := template.New("test").Funcs(map[string]interface{}{
	//	"echo": func(params ...interface{}) interface{} {
	//		return fmt.Sprintf("echo: %v", params[0])
	//	},
	//})
	//t, err := tpl.Parse("{{lt .age 20 | echo}}")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//var buf = &bytes.Buffer{}
	//err = t.Execute(buf, map[string]interface{}{
	//	"age": 19,
	//}) // 执行模板解析
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(buf.String())

	fmt.Println(goft.ExecExpr("myage < 20", map[string]interface{}{"myage": 19}))
}
