package main

import (
	"fmt"
	"gin-up/BeanExpr/FuncExpr"
	parser "gin-up/test"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"log"
	"reflect"
	"strconv"
	"strings"
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

type FuncExprListener struct {
	*FuncExpr.BaseBeanExprListener
	funcName string
	args     []reflect.Value
}

// ExitMethodCall is called when production methodCall is exited.
func (this *FuncExprListener) ExitMethodCall(ctx *FuncExpr.MethodCallContext) {
	log.Println("method 内容是： ", ctx.GetText())
}

// ExitFuncCall is called when production FuncCall is exited.
func (this *FuncExprListener) ExitFuncCall(ctx *FuncExpr.FuncCallContext) {
	log.Println("函数名是: ", ctx.GetStart().GetText())
	this.funcName = ctx.GetStart().GetText()
}

// EnterFuncArgs is called when production FuncArgs is entered.
func (this *FuncExprListener) EnterFuncArgs(ctx *FuncExpr.FuncArgsContext) {
	for i := 0; i < ctx.GetChildCount(); i++ {
		token := ctx.GetChild(i).GetPayload().(*antlr.CommonToken)
		var value reflect.Value
		switch token.GetTokenType() {
		case FuncExpr.BeanExprLexerStringArg:
			stringArg := strings.Trim(token.GetText(), "'")
			value = reflect.ValueOf(stringArg)
			break
		case FuncExpr.BeanExprLexerIntArg:
			v, err := strconv.ParseInt(token.GetText(), 10, 64)
			if err != nil {
				panic("parse int64 error")
			}
			value = reflect.ValueOf(v)
			break
		case FuncExpr.BeanExprLexerFloatArg:
			v, err := strconv.ParseFloat(token.GetText(), 64)
			if err != nil {
				panic("parse float64 error")
			}
			value = reflect.ValueOf(v)
			break
		default:
			continue
		}
		this.args = append(this.args, value)
	}
}

func (this *FuncExprListener) Run() {
	if f, ok := FuncMap[this.funcName]; ok {
		v := reflect.ValueOf(f)
		if v.Kind() == reflect.Func {
			v.Call(this.args)
		}
	}
}

var FuncMap map[string]interface{}

func main() {
	FuncMap = map[string]interface{}{
		"test": func(name string, age int64) {
			log.Println("this is ", name, " and age is: ", age)
		},
	}

	is := antlr.NewInputStream("User.test('custer',16)")
	lexer := FuncExpr.NewBeanExprLexer(is)
	ts := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := FuncExpr.NewBeanExprParser(ts)
	lis := &FuncExprListener{}
	antlr.ParseTreeWalkerDefault.Walk(lis, p.Start())
	lis.Run()
}

func main2() {
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
