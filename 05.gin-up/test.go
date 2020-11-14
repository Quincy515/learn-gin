package main

import (
	"fmt"
	parser "gin-up/test"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type Me struct{}

func (this *Me) Age() int {
	return 21
}
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

	//fmt.Println(goft.ExecExpr("myage < 20", map[string]interface{}{"myage": 19}))
	//fmt.Println(goft.ExecExpr("gt .me.Age .myage",
	//	map[string]interface{}{"myage": 19, "me": &Me{}}))
	is := antlr.NewInputStream("1+2*3")
	lexer := parser.NewCalcLexer(is) // 获取语法解析器对象
	for {
		t := lexer.NextToken()
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}
		fmt.Printf("%s (%q)\n",
			lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	}
}
