%{
package parsing

import(
	"github.com/paulfchristiano/dwimmer/term"
	"github.com/paulfchristiano/dwimmer/data/represent"
)
%}

%token <string> WHITE PROSE WORD SYMBOL MAKE
%token <int> NUM
%token <string> '[' ']' '(' ')' '{' '}' ',' ':' '|' '"' '?' '!' '-' '.' '#' '@'
%token WANT_ACTION WANT_TERM

%type <expr> expr clauselist exprbracket
%type <clause> clause int str word
%type <string> text textblock quotedtext

%union{
    string string
    int int
    rune rune
    clause ExprPart
    expr *Expr
}

%%

result: 
      WANT_ACTION WORD WHITE NUM WHITE expr {yylex.(*lexer).setActionResult($2, $6, $4)}
    | WANT_ACTION WORD WHITE expr {yylex.(*lexer).setActionResult($2, $4, -1)}
    | WANT_ACTION WORD optionalwhitespace '@' optionalwhitespace WORD WHITE expr {yylex.(*lexer).setTransitiveActionResult($2, $6, $8)}
    | WANT_ACTION WORD {yylex.(*lexer).setActionResult($2, nil, -1)}
    | WANT_TERM expr {yylex.(*lexer).setTermResult($2)}
    | WANT_ACTION WORD WHITE NUM {{yylex.(*lexer).setActionResult($2, nil, $4)}}

expr: clauselist

clauselist: {$$ = EmptyExpr()}
    | clauselist clause { $$ = $1; $$.append($2) }

clause: int | str | word 
    | text {$$ = exprText($1)}
    | exprbracket {$$ = exprTerm{toC($1)}}

exprbracket: '[' expr ']' {$$ = $2}
    | '(' expr ')' {$$ = $2}

str: '"' quotedtext '"' {$$ = exprTerm{term.ConstC{represent.Str($2)}}}
   | '"' '"' {$$ = exprTerm{term.ConstC{represent.Str("")}}}

int: NUM { $$ = exprTerm{term.ConstC{represent.Int($1)}} }

word: WORD {$$ = yylex.(*lexer).parseWord($1)}

quotedtext: text | WORD | quotedtext quotedtext {$$ = $1+$2}

optionalwhitespace : | WHITE

textblock: WHITE | PROSE | SYMBOL
         | ',' | ':' | '?' | '!' | '-' | '.'

text: textblock
    | textblock text {$$ = $1 + $2}

%%
