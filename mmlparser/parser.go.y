%{
package mmlparser
%}

%union {
    node        *Node
    token       Token
    str         string
}

// プログラムの構成要素を指定
%type<node> program block line tone loop expr
%type<token> toneName toneFlag
%type<str> toneFlags
// トークンの定義
%token<token> LF WORD NUMBER
%token<token> 'c' 'd' 'e' 'f' 'g' 'a' 'b' '#' '+' '-' '*' '[' ']' ':'
%%

// 文法規則を指定
program
    : block             { $$ = $1; yylex.(*Lexer).result = $$ }
 
block
    :                   { $$ = NewNode(Nop) }
    | line              { $$ = NewNode(NodeList); $$ = AppendChildNode($$, $1) }
    | block line        { $$ = AppendChildNode($1, $2) }

line
    : tone
    | loop
    | LF                { $$ = NewNode(NodeEOL) }


expr
    : NUMBER            { $$ = NewNumberNode($1) }

loop
    : '[' expr          { $$ = NewLoopNodeBegin($1, $2) }
    | '['               { $$ = NewLoopNodeBegin($1, nil) }
    | ']'               { $$ = NewLoopNodeEnd($1)   }
    | ':'               { $$ = NewLoopNodeBreak($1) }

tone
    : toneName toneFlags        { $$ = NewToneNode($1, $2) }
    | toneName                  { $$ = NewToneNode($1, "") }

toneName
    : 'c' | 'd' | 'e' | 'f' | 'g' | 'a' | 'b'

toneFlags
    : toneFlag              { $$ = $1.label      }
    | toneFlags toneFlag    { $$ = $1 + $2.label }

toneFlag
    : '#' | '+' | '-'


%%

