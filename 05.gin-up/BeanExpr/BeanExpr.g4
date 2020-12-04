grammar BeanExpr;

// Rules
start : functionCall EOF;

functionCall
    : FuncName '(' functionArgs? ')'  #FuncCall   //函数调用  譬如 GetAge()  或 GetUser(101)
    ;

//参数目前先支持 字符串或 数字
functionArgs
            : (FloatArg | StringArg | IntArg) (',' (FloatArg | StringArg | IntArg))* #FuncArgs //带一个，为分隔符的序列
            ;
//fragment标识的规则只能为其它词法规则提供基础
fragment DIGIT: [0-9];

//以下是Token
//字符串 参数 用单引号  如 'abc' 也可以是 '\'abc\'' 支持单引号转义
StringArg: '\'' ('\\'. | '\'\'' | ~('\'' | '\\'))* '\'';
FloatArg: '-'?Float;
IntArg: '-'?DIGIT+;
FuncName: [a-zA-Z][a-zA-Z0-9]*; //函数名称 必须字母开头, 支持数字字母的组合
Dot: '.';

Float :(DIGIT+)? '.' DIGIT+   // 如 19.02  .02   目前不支持科学计数
       ;
WHITESPACE: [ \r\n\t]+ -> skip;
