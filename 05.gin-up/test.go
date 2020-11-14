package main

import (
	"fmt"
	parser "gin-up/test"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"strconv"
)

type Me struct{}

func (this *Me) Age() int {
	return 21
}

type calcListener struct {
	*parser.BaseCalcListener
	nums []int // 保存数据
}

func (this *calcListener) push(i int) {
	this.nums = append(this.nums, i)
}
func (this *calcListener) pop() int {
	if len(this.nums) < 1 {
		panic("unable to pop")
	}
	result := this.nums[len(this.nums)-1]
	this.nums = this.nums[:len(this.nums)-1]
	return result
}
func (this *calcListener) ExitAddSub(ctx *parser.AddSubContext) {
	n1, n2 := this.pop(), this.pop()
	switch ctx.GetOp().GetTokenType() {
	case parser.CalcParserADD:
		this.push(n1 + n2)
		break
	case parser.CalcLexerSUB:
		this.push(n2 - n1)
		break
	}
}
func (this *calcListener) ExitNumber(ctx *parser.NumberContext) {
	num, _ := strconv.Atoi(ctx.GetText())
	this.nums = append(this.nums, num)
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
	is := antlr.NewInputStream("9-2")
	lexer := parser.NewCalcLexer(is) // 获取语法解析器对象
	//for {
	//	t := lexer.NextToken()
	//	if t.GetTokenType() == antlr.TokenEOF {
	//		break
	//	}
	//	fmt.Printf("%s (%q)\n",
	//		lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	//}
	ts := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewCalcParser(ts)
	lis := &calcListener{}
	antlr.ParseTreeWalkerDefault.Walk(lis, p.Start())
	fmt.Println(lis.pop())
}
