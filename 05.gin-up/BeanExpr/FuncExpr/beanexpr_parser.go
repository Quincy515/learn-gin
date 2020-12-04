// Code generated from /Users/tianxiaoqiang/Work/2020/study/learn-gin/05.gin-up/BeanExpr/BeanExpr.g4 by ANTLR 4.8. DO NOT EDIT.

package FuncExpr // BeanExpr
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa


var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 12, 27, 4, 
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 5, 
	3, 15, 10, 3, 3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 7, 4, 22, 10, 4, 12, 4, 14, 
	4, 25, 11, 4, 3, 4, 2, 2, 5, 2, 4, 6, 2, 2, 2, 25, 2, 8, 3, 2, 2, 2, 4, 
	11, 3, 2, 2, 2, 6, 18, 3, 2, 2, 2, 8, 9, 5, 4, 3, 2, 9, 10, 7, 2, 2, 3, 
	10, 3, 3, 2, 2, 2, 11, 12, 7, 7, 2, 2, 12, 14, 7, 3, 2, 2, 13, 15, 5, 6, 
	4, 2, 14, 13, 3, 2, 2, 2, 14, 15, 3, 2, 2, 2, 15, 16, 3, 2, 2, 2, 16, 17, 
	7, 4, 2, 2, 17, 5, 3, 2, 2, 2, 18, 23, 7, 9, 2, 2, 19, 20, 7, 5, 2, 2, 
	20, 22, 7, 9, 2, 2, 21, 19, 3, 2, 2, 2, 22, 25, 3, 2, 2, 2, 23, 21, 3, 
	2, 2, 2, 23, 24, 3, 2, 2, 2, 24, 7, 3, 2, 2, 2, 25, 23, 3, 2, 2, 2, 4, 
	14, 23,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'('", "')'", "','", "", "", "'.'",
}
var symbolicNames = []string{
	"", "", "", "", "StringArg", "FuncName", "Dot", "Args", "Int", "Number", 
	"WHITESPACE",
}

var ruleNames = []string{
	"start", "functionCall", "functionArgs",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type BeanExprParser struct {
	*antlr.BaseParser
}

func NewBeanExprParser(input antlr.TokenStream) *BeanExprParser {
	this := new(BeanExprParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "BeanExpr.g4"

	return this
}

// BeanExprParser tokens.
const (
	BeanExprParserEOF = antlr.TokenEOF
	BeanExprParserT__0 = 1
	BeanExprParserT__1 = 2
	BeanExprParserT__2 = 3
	BeanExprParserStringArg = 4
	BeanExprParserFuncName = 5
	BeanExprParserDot = 6
	BeanExprParserArgs = 7
	BeanExprParserInt = 8
	BeanExprParserNumber = 9
	BeanExprParserWHITESPACE = 10
)

// BeanExprParser rules.
const (
	BeanExprParserRULE_start = 0
	BeanExprParserRULE_functionCall = 1
	BeanExprParserRULE_functionArgs = 2
)

// IStartContext is an interface to support dynamic dispatch.
type IStartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = BeanExprParserRULE_start
	return p
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = BeanExprParserRULE_start

	return p
}

func (s *StartContext) GetParser() antlr.Parser { return s.parser }

func (s *StartContext) FunctionCall() IFunctionCallContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *StartContext) EOF() antlr.TerminalNode {
	return s.GetToken(BeanExprParserEOF, 0)
}

func (s *StartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.EnterStart(s)
	}
}

func (s *StartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.ExitStart(s)
	}
}




func (p *BeanExprParser) Start() (localctx IStartContext) {
	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, BeanExprParserRULE_start)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(6)
		p.FunctionCall()
	}
	{
		p.SetState(7)
		p.Match(BeanExprParserEOF)
	}



	return localctx
}


// IFunctionCallContext is an interface to support dynamic dispatch.
type IFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunctionCallContext differentiates from other interfaces.
	IsFunctionCallContext()
}

type FunctionCallContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallContext() *FunctionCallContext {
	var p = new(FunctionCallContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = BeanExprParserRULE_functionCall
	return p
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = BeanExprParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) CopyFrom(ctx *FunctionCallContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}




type FuncCallContext struct {
	*FunctionCallContext
}

func NewFuncCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FuncCallContext {
	var p = new(FuncCallContext)

	p.FunctionCallContext = NewEmptyFunctionCallContext()
	p.parser = parser
	p.CopyFrom(ctx.(*FunctionCallContext))

	return p
}

func (s *FuncCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncCallContext) FuncName() antlr.TerminalNode {
	return s.GetToken(BeanExprParserFuncName, 0)
}

func (s *FuncCallContext) FunctionArgs() IFunctionArgsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionArgsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionArgsContext)
}


func (s *FuncCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.EnterFuncCall(s)
	}
}

func (s *FuncCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.ExitFuncCall(s)
	}
}



func (p *BeanExprParser) FunctionCall() (localctx IFunctionCallContext) {
	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, BeanExprParserRULE_functionCall)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	localctx = NewFuncCallContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(9)
		p.Match(BeanExprParserFuncName)
	}
	{
		p.SetState(10)
		p.Match(BeanExprParserT__0)
	}
	p.SetState(12)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == BeanExprParserArgs {
		{
			p.SetState(11)
			p.FunctionArgs()
		}

	}
	{
		p.SetState(14)
		p.Match(BeanExprParserT__1)
	}



	return localctx
}


// IFunctionArgsContext is an interface to support dynamic dispatch.
type IFunctionArgsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunctionArgsContext differentiates from other interfaces.
	IsFunctionArgsContext()
}

type FunctionArgsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionArgsContext() *FunctionArgsContext {
	var p = new(FunctionArgsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = BeanExprParserRULE_functionArgs
	return p
}

func (*FunctionArgsContext) IsFunctionArgsContext() {}

func NewFunctionArgsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionArgsContext {
	var p = new(FunctionArgsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = BeanExprParserRULE_functionArgs

	return p
}

func (s *FunctionArgsContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionArgsContext) CopyFrom(ctx *FunctionArgsContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *FunctionArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionArgsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}




type FuncArgsContext struct {
	*FunctionArgsContext
}

func NewFuncArgsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FuncArgsContext {
	var p = new(FuncArgsContext)

	p.FunctionArgsContext = NewEmptyFunctionArgsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*FunctionArgsContext))

	return p
}

func (s *FuncArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncArgsContext) AllArgs() []antlr.TerminalNode {
	return s.GetTokens(BeanExprParserArgs)
}

func (s *FuncArgsContext) Args(i int) antlr.TerminalNode {
	return s.GetToken(BeanExprParserArgs, i)
}


func (s *FuncArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.EnterFuncArgs(s)
	}
}

func (s *FuncArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.ExitFuncArgs(s)
	}
}



func (p *BeanExprParser) FunctionArgs() (localctx IFunctionArgsContext) {
	localctx = NewFunctionArgsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, BeanExprParserRULE_functionArgs)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	localctx = NewFuncArgsContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(16)
		p.Match(BeanExprParserArgs)
	}
	p.SetState(21)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == BeanExprParserT__2 {
		{
			p.SetState(17)
			p.Match(BeanExprParserT__2)
		}
		{
			p.SetState(18)
			p.Match(BeanExprParserArgs)
		}


		p.SetState(23)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


