// Code generated from /Users/tianxiaoqiang/Work/2020/study/learn-gin/05.gin-up/BeanExpr/BeanExpr.g4 by ANTLR 4.8. DO NOT EDIT.

package FuncExpr // BeanExpr
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BeanExprListener is a complete listener for a parse tree produced by BeanExprParser.
type BeanExprListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterFuncCall is called when entering the FuncCall production.
	EnterFuncCall(c *FuncCallContext)

	// EnterFuncArgs is called when entering the FuncArgs production.
	EnterFuncArgs(c *FuncArgsContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitFuncCall is called when exiting the FuncCall production.
	ExitFuncCall(c *FuncCallContext)

	// ExitFuncArgs is called when exiting the FuncArgs production.
	ExitFuncArgs(c *FuncArgsContext)
}
