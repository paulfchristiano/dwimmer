//line parser.y:2
package parsing

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/term"
)

//line parser.y:19
type yySymType struct {
	yys    int
	string string
	int    int
	rune   rune
	clause ExprPart
	expr   *Expr
}

const WHITE = 57346
const PROSE = 57347
const WORD = 57348
const SYMBOL = 57349
const MAKE = 57350
const NUM = 57351
const WANT_ACTION = 57352
const WANT_TERM = 57353

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"WHITE",
	"PROSE",
	"WORD",
	"SYMBOL",
	"MAKE",
	"NUM",
	"'['",
	"']'",
	"'('",
	"')'",
	"'{'",
	"'}'",
	"','",
	"':'",
	"'|'",
	"'\"'",
	"'?'",
	"'!'",
	"'-'",
	"'.'",
	"'#'",
	"'@'",
	"WANT_ACTION",
	"WANT_TERM",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:66
//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 4,
	25, 24,
	-2, 4,
	-1, 7,
	25, 25,
	-2, 8,
}

const yyNprod = 37
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 99

var yyAct = [...]int{

	5, 8, 2, 3, 32, 46, 45, 30, 31, 48,
	21, 22, 17, 23, 4, 15, 19, 35, 20, 49,
	38, 39, 24, 25, 13, 16, 26, 27, 28, 29,
	21, 22, 36, 23, 41, 42, 37, 40, 44, 7,
	1, 47, 24, 25, 18, 43, 26, 27, 28, 29,
	50, 21, 22, 36, 23, 33, 12, 11, 10, 21,
	22, 36, 23, 24, 25, 9, 34, 26, 27, 28,
	29, 24, 25, 14, 6, 26, 27, 28, 29, 21,
	22, 0, 23, 0, 0, 0, 0, 0, 0, 0,
	0, 24, 25, 0, 0, 26, 27, 28, 29,
}
var yyPact = [...]int{

	-24, -1000, 8, -1000, 35, -1000, 6, -2, -21, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 47, -1000, 75, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	33, -1000, 31, 26, -1000, -1000, -1000, -1000, -5, -8,
	-1000, 3, -1000, -1000, 55, -1000, -1000, -1000, 15, -1000,
	-1000,
}
var yyPgo = [...]int{

	0, 0, 74, 73, 65, 58, 57, 56, 17, 44,
	38, 40, 1,
}
var yyR1 = [...]int{

	0, 11, 11, 11, 11, 11, 11, 1, 2, 2,
	4, 4, 4, 4, 4, 3, 3, 6, 6, 5,
	7, 10, 10, 10, 12, 12, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 8, 8,
}
var yyR2 = [...]int{

	0, 6, 4, 8, 2, 2, 4, 1, 0, 2,
	1, 1, 1, 1, 1, 3, 3, 3, 2, 1,
	1, 1, 1, 2, 0, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 2,
}
var yyChk = [...]int{

	-1000, -11, 26, 27, 6, -1, -2, 4, -12, -4,
	-5, -6, -7, -8, -3, 9, 19, 6, -9, 10,
	12, 4, 5, 7, 16, 17, 20, 21, 22, 23,
	9, -1, 25, -10, 19, -8, 6, -8, -1, -1,
	4, -12, 4, 19, -10, 11, 13, -1, 6, 4,
	-1,
}
var yyDef = [...]int{

	0, -2, 0, 8, -2, 5, 7, -2, 0, 9,
	10, 11, 12, 13, 14, 19, 0, 20, 35, 8,
	8, 26, 27, 28, 29, 30, 31, 32, 33, 34,
	6, 2, 24, 0, 18, 21, 22, 36, 0, 0,
	8, 0, 25, 17, 23, 15, 16, 1, 0, 8,
	3,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 21, 19, 24, 3, 3, 3, 3,
	12, 13, 3, 3, 16, 22, 23, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 17, 3,
	3, 3, 3, 20, 25, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 10, 3, 11, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 14, 18, 15,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 26, 27,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lookahead func() int
}

func (p *yyParserImpl) Lookahead() int {
	return p.lookahead()
}

func yyNewParser() yyParser {
	p := &yyParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yytoken := -1 // yychar translated into internal numbering
	yyrcvr.lookahead = func() int { return yychar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yychar = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar, yytoken = yylex1(yylex, &yylval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yychar = -1
		yytoken = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar, yytoken = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yychar = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:30
		{
			yylex.(*lexer).setActionResult(yyDollar[2].string, yyDollar[6].expr, yyDollar[4].int)
		}
	case 2:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:31
		{
			yylex.(*lexer).setActionResult(yyDollar[2].string, yyDollar[4].expr, -1)
		}
	case 3:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:32
		{
			yylex.(*lexer).setTransitiveActionResult(yyDollar[2].string, yyDollar[6].string, yyDollar[8].expr)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:33
		{
			yylex.(*lexer).setActionResult(yyDollar[2].string, nil, -1)
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:34
		{
			yylex.(*lexer).setTermResult(yyDollar[2].expr)
		}
	case 6:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:35
		{
			{
				yylex.(*lexer).setActionResult(yyDollar[2].string, nil, yyDollar[4].int)
			}
		}
	case 8:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:39
		{
			yyVAL.expr = EmptyExpr()
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:40
		{
			yyVAL.expr = yyDollar[1].expr
			yyVAL.expr.append(yyDollar[2].clause)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:43
		{
			yyVAL.clause = exprText(yyDollar[1].string)
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:44
		{
			yyVAL.clause = exprTerm{toC(yyDollar[1].expr)}
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:46
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:47
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:49
		{
			yyVAL.clause = exprTerm{term.ConstC{represent.Str(yyDollar[2].string)}}
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:50
		{
			yyVAL.clause = exprTerm{term.ConstC{represent.Str("")}}
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:52
		{
			yyVAL.clause = exprTerm{term.ConstC{represent.Int(yyDollar[1].int)}}
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:54
		{
			yyVAL.clause = yylex.(*lexer).parseWord(yyDollar[1].string)
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:56
		{
			yyVAL.string = yyDollar[1].string + yyDollar[2].string
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:64
		{
			yyVAL.string = yyDollar[1].string + yyDollar[2].string
		}
	}
	goto yystack /* stack new state and value */
}
