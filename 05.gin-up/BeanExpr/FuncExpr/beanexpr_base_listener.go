// Code generated from /Users/tianxiaoqiang/Work/2020/study/learn-gin/05.gin-up/BeanExpr/BeanExpr.g4 by ANTLR 4.8. DO NOT EDIT.

package FuncExpr // BeanExpr
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseBeanExprListener is a complete listener for a parse tree produced by BeanExprParser.
type BaseBeanExprListener struct{}

var _ BeanExprListener = &BaseBeanExprListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseBeanExprListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseBeanExprListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseBeanExprListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseBeanExprListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseBeanExprListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseBeanExprListener) ExitStart(ctx *StartContext) {}

// EnterFuncCall is called when production FuncCall is entered.
func (s *BaseBeanExprListener) EnterFuncCall(ctx *FuncCallContext) {}

// ExitFuncCall is called when production FuncCall is exited.
func (s *BaseBeanExprListener) ExitFuncCall(ctx *FuncCallContext) {}

// EnterFuncArgs is called when production FuncArgs is entered.
func (s *BaseBeanExprListener) EnterFuncArgs(ctx *FuncArgsContext) {}

// ExitFuncArgs is called when production FuncArgs is exited.
func (s *BaseBeanExprListener) ExitFuncArgs(ctx *FuncArgsContext) {}
