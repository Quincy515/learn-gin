grammar Calc;

// Tokens
MUL: '*';
DIV: '/';
ADD: '+';
SUB: '-';
NUMBER: [0-9]+;
WHITESPACE: [ \r\n\t]+ -> skip;

// Rules
start : expr EOF;

expr
   : expr op=('*'|'/') expr # MulDiv
   | expr op=('+'|'-') expr # AddSub
   | NUMBER                             # Number
   ;