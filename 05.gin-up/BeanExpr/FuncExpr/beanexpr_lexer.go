// Code generated from /Users/tianxiaoqiang/Work/2020/study/learn-gin/05.gin-up/BeanExpr/BeanExpr.g4 by ANTLR 4.8. DO NOT EDIT.

package FuncExpr
import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)
// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter


var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 12, 94, 8, 
	1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 
	7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 3, 
	2, 3, 2, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 
	6, 3, 6, 7, 6, 40, 10, 6, 12, 6, 14, 6, 43, 11, 6, 3, 6, 3, 6, 3, 7, 3, 
	7, 7, 7, 49, 10, 7, 12, 7, 14, 7, 52, 11, 7, 3, 8, 3, 8, 3, 9, 5, 9, 57, 
	10, 9, 3, 9, 3, 9, 5, 9, 61, 10, 9, 3, 10, 3, 10, 7, 10, 65, 10, 10, 12, 
	10, 14, 10, 68, 11, 10, 3, 11, 5, 11, 71, 10, 11, 3, 11, 3, 11, 6, 11, 
	75, 10, 11, 13, 11, 14, 11, 76, 3, 11, 3, 11, 7, 11, 81, 10, 11, 12, 11, 
	14, 11, 84, 11, 11, 5, 11, 86, 10, 11, 3, 12, 6, 12, 89, 10, 12, 13, 12, 
	14, 12, 90, 3, 12, 3, 12, 2, 2, 13, 3, 3, 5, 4, 7, 5, 9, 2, 11, 6, 13, 
	7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 3, 2, 8, 3, 2, 50, 59, 4, 2, 41, 
	41, 94, 94, 4, 2, 67, 92, 99, 124, 5, 2, 50, 59, 67, 92, 99, 124, 3, 2, 
	51, 59, 5, 2, 11, 12, 15, 15, 34, 34, 2, 104, 2, 3, 3, 2, 2, 2, 2, 5, 3, 
	2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 
	3, 2, 2, 2, 2, 17, 3, 2, 2, 2, 2, 19, 3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 
	23, 3, 2, 2, 2, 3, 25, 3, 2, 2, 2, 5, 27, 3, 2, 2, 2, 7, 29, 3, 2, 2, 2, 
	9, 31, 3, 2, 2, 2, 11, 33, 3, 2, 2, 2, 13, 46, 3, 2, 2, 2, 15, 53, 3, 2, 
	2, 2, 17, 60, 3, 2, 2, 2, 19, 62, 3, 2, 2, 2, 21, 85, 3, 2, 2, 2, 23, 88, 
	3, 2, 2, 2, 25, 26, 7, 42, 2, 2, 26, 4, 3, 2, 2, 2, 27, 28, 7, 43, 2, 2, 
	28, 6, 3, 2, 2, 2, 29, 30, 7, 46, 2, 2, 30, 8, 3, 2, 2, 2, 31, 32, 9, 2, 
	2, 2, 32, 10, 3, 2, 2, 2, 33, 41, 7, 41, 2, 2, 34, 35, 7, 94, 2, 2, 35, 
	40, 11, 2, 2, 2, 36, 37, 7, 41, 2, 2, 37, 40, 7, 41, 2, 2, 38, 40, 10, 
	3, 2, 2, 39, 34, 3, 2, 2, 2, 39, 36, 3, 2, 2, 2, 39, 38, 3, 2, 2, 2, 40, 
	43, 3, 2, 2, 2, 41, 39, 3, 2, 2, 2, 41, 42, 3, 2, 2, 2, 42, 44, 3, 2, 2, 
	2, 43, 41, 3, 2, 2, 2, 44, 45, 7, 41, 2, 2, 45, 12, 3, 2, 2, 2, 46, 50, 
	9, 4, 2, 2, 47, 49, 9, 5, 2, 2, 48, 47, 3, 2, 2, 2, 49, 52, 3, 2, 2, 2, 
	50, 48, 3, 2, 2, 2, 50, 51, 3, 2, 2, 2, 51, 14, 3, 2, 2, 2, 52, 50, 3, 
	2, 2, 2, 53, 54, 7, 48, 2, 2, 54, 16, 3, 2, 2, 2, 55, 57, 7, 47, 2, 2, 
	56, 55, 3, 2, 2, 2, 56, 57, 3, 2, 2, 2, 57, 58, 3, 2, 2, 2, 58, 61, 5, 
	21, 11, 2, 59, 61, 5, 11, 6, 2, 60, 56, 3, 2, 2, 2, 60, 59, 3, 2, 2, 2, 
	61, 18, 3, 2, 2, 2, 62, 66, 9, 6, 2, 2, 63, 65, 5, 9, 5, 2, 64, 63, 3, 
	2, 2, 2, 65, 68, 3, 2, 2, 2, 66, 64, 3, 2, 2, 2, 66, 67, 3, 2, 2, 2, 67, 
	20, 3, 2, 2, 2, 68, 66, 3, 2, 2, 2, 69, 71, 5, 19, 10, 2, 70, 69, 3, 2, 
	2, 2, 70, 71, 3, 2, 2, 2, 71, 72, 3, 2, 2, 2, 72, 74, 7, 48, 2, 2, 73, 
	75, 5, 9, 5, 2, 74, 73, 3, 2, 2, 2, 75, 76, 3, 2, 2, 2, 76, 74, 3, 2, 2, 
	2, 76, 77, 3, 2, 2, 2, 77, 86, 3, 2, 2, 2, 78, 82, 9, 6, 2, 2, 79, 81, 
	5, 9, 5, 2, 80, 79, 3, 2, 2, 2, 81, 84, 3, 2, 2, 2, 82, 80, 3, 2, 2, 2, 
	82, 83, 3, 2, 2, 2, 83, 86, 3, 2, 2, 2, 84, 82, 3, 2, 2, 2, 85, 70, 3, 
	2, 2, 2, 85, 78, 3, 2, 2, 2, 86, 22, 3, 2, 2, 2, 87, 89, 9, 7, 2, 2, 88, 
	87, 3, 2, 2, 2, 89, 90, 3, 2, 2, 2, 90, 88, 3, 2, 2, 2, 90, 91, 3, 2, 2, 
	2, 91, 92, 3, 2, 2, 2, 92, 93, 8, 12, 2, 2, 93, 24, 3, 2, 2, 2, 14, 2, 
	39, 41, 50, 56, 60, 66, 70, 76, 82, 85, 90, 3, 8, 2, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'('", "')'", "','", "", "", "'.'",
}

var lexerSymbolicNames = []string{
	"", "", "", "", "StringArg", "FuncName", "Dot", "Args", "Int", "Number", 
	"WHITESPACE",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "T__2", "DIGIT", "StringArg", "FuncName", "Dot", "Args", 
	"Int", "Number", "WHITESPACE",
}

type BeanExprLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewBeanExprLexer(input antlr.CharStream) *BeanExprLexer {

	l := new(BeanExprLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "BeanExpr.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// BeanExprLexer tokens.
const (
	BeanExprLexerT__0 = 1
	BeanExprLexerT__1 = 2
	BeanExprLexerT__2 = 3
	BeanExprLexerStringArg = 4
	BeanExprLexerFuncName = 5
	BeanExprLexerDot = 6
	BeanExprLexerArgs = 7
	BeanExprLexerInt = 8
	BeanExprLexerNumber = 9
	BeanExprLexerWHITESPACE = 10
)

