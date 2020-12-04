package main

import (
	"gin-up/BeanExpr/FuncExpr"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type FuncExprListener_ struct {
	*FuncExpr.BaseBeanExprListener
	funcName   string
	args       []reflect.Value
	methodName string // 方法名 user.abc.bcd.GetAge()
	execType   uint8  // 执行类型  0代表函数 也就是默认值， 1代表struct执行
}

func (this *FuncExprListener_) ExitMethodCall(ctx *FuncExpr.MethodCallContext) {
	this.execType = 1
	this.methodName = ctx.GetStart().GetText()

}
func (this *FuncExprListener_) ExitFuncCall(ctx *FuncExpr.FuncCallContext) {
	//log.Println("函数名是:",	ctx.GetStart().GetText())
	this.funcName = ctx.GetStart().GetText()
}
func (this *FuncExprListener_) ExitFuncArgs(ctx *FuncExpr.FuncArgsContext) {
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
func (this *FuncExprListener_) findField(method string, v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if field := v.FieldByName(method); field.IsValid() {
		return field
	}
	return reflect.Value{}
}
func (this *FuncExprListener_) Run() {
	switch this.execType {
	case 0:
		if f, ok := FuncMap_[this.funcName]; ok {
			v := reflect.ValueOf(f)
			if v.Kind() == reflect.Func {
				v.Call(this.args)
			}
		}
		break
	case 1: // struct方法执行
		ms := strings.Split(this.methodName, ".")
		if obj, ok := FuncMap_[ms[0]]; ok {
			objv := reflect.ValueOf(obj)
			current := objv
			for i := 1; i < len(ms); i++ {
				if i == len(ms)-1 { //最后一个是方法名
					if method := current.MethodByName(ms[i]); !method.IsValid() {
						panic("method error:" + ms[i])
					} else {
						method.Call(this.args)
					}
					break
				}
				field := this.findField(ms[i], current)
				if field.IsValid() {
					current = field
				} else {
					panic("field error:" + ms[i])
				}
			}
		}
	default:
		log.Println("nothing to do")
	}
}

type Admin struct{}

func (this *Admin) Abc() {
	log.Println("admin")
}

type User struct {
	Adm *Admin
}

func (this *User) Name(prefix string, age int64) {
	log.Println(prefix+"custer", age)
}

var FuncMap_ map[string]interface{}

func main() {
	FuncMap_ = map[string]interface{}{
		"test": func(name string, age int64) {
			log.Println("this is ", name, " and  age is :", age)
		},
		"User": &User{Adm: &Admin{}},
	}
	//is := antlr.NewInputStream("User.Adm.Abc()")
	is := antlr.NewInputStream("User.Name('test',19)")

	lexer := FuncExpr.NewBeanExprLexer(is)
	ts := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := FuncExpr.NewBeanExprParser(ts)
	lis := &FuncExprListener_{}
	antlr.ParseTreeWalkerDefault.Walk(lis, p.Start())

	lis.Run()
}
