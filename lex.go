package sakuramml

import (
	"fmt"
	"strconv"
	"strings"
)

// Token を表す
type Token struct {
	label    string
	value    float64
	lineinfo LineInfo
}

var tokenName = map[int]string{
	LF:     "LF",
	WORD:   "WORD",
	NUMBER: "NUMBER",
}

type Lexer struct {
	slexer     SLexer
	result     *Node
	parseError error
	lastToken  *Token
}

// Parse : コードを実行
func Parse(code string, fileno int) (*Node, error) {
	lexer := Lexer{
		slexer:     *newSLexer(code+"\n", fileno),
		result:     nil,
		lastToken:  nil,
		parseError: nil,
	}
	yyDebug = 1
	yyErrorVerbose = true
	yyParse(&lexer)
	if lexer.parseError != nil {
		return nil, lexer.parseError
	}
	return lexer.result, nil
}

// Lex : トークンを一つずつ返す
func (p *Lexer) Lex(lval *yySymType) int {
	tok := p.getToken(lval)
	p.lastToken = &lval.token
	SakuraLog(fmt.Sprintf("[Lex] %s", getTokenName(tok)))
	return tok
}

func getTokenName(tok int) string {
	if 0x21 <= tok && tok <= 0x7e {
		return string(rune(tok))
	}
	w, ok := tokenName[tok]
	if ok {
		return w
	}
	return fmt.Sprintf("tok:%d", tok)
}

func (p *Lexer) getToken(lval *yySymType) int {
	p.slexer.skipSpace()
	// Set LineInfo
	CurLine.LineNo = p.slexer.lineno
	CurLine.FileNo = p.slexer.fileno
	// Check EOF
	if p.slexer.isEOF() {
		return 0 // end
	}
	c := p.slexer.peek()
	// LF ?
	if c == '\n' {
		lval.token = p.newToken("\n")
		p.slexer.lineno++
		p.slexer.next()
		return LF
	}
	// Comment ?
	if p.slexer.testStr("/*") {
		p.slexer.index += 2
		p.slexer.lineno += p.slexer.skipTo("*/")
		return COMMENT
	}
	// NUMBER ?
	if isDigit(c) || c == '^' || c == '%' || (c == '-' && isDigit(p.slexer.peekNext())) {
		return p.lexNumber(lval)
	}
	// 演算子 ?
	if isFlag(c) {
		p.slexer.next()
		lval.token = p.newToken(string(c))
		return int(c)
	}
	// 小文字コマンド
	if 'a' <= c && c <= 'z' {
		lval.token = p.newToken(string(c))
		p.slexer.next()
		return int(c)
	}
	// 大文字コマンド
	if 'A' <= c && c <= 'Z' || c == '_' {
		return p.lexWord(lval)
	}
	// 記号
	if c == '{' {
		p.slexer.next()
		lval.token = p.newToken("{")
		return PAREN_L
	}
	if c == '}' {
		p.slexer.next()
		lval.token = p.newToken("}")
		return PAREN_R
	}
	if ('!' <= c && c <= '/') || (':' <= c && c <= '@') || ('[' <= c && c <= '`') || ('{' <= c && c <= '~') {
		p.slexer.next()
		lval.token = p.newToken(string(c))
		return int(c)
	}
	return -1
}

func (p *Lexer) lexNumber(lval *yySymType) int {
	s := ""
	if p.slexer.peek() == '-' {
		s += "-"
		p.slexer.next()
	}
	for !p.slexer.isEOF() {
		c := p.slexer.peek()
		if isDigit(c) || c == '.' || c == '^' || c == '%' {
			s += string(c)
			p.slexer.next()
			continue
		}
		p.slexer.skipSpace()
		break
	}
	lval.token = p.newToken(s)
	lval.token.value, _ = strconv.ParseFloat(s, 64)
	return NUMBER
}

func (p *Lexer) lexWord(lval *yySymType) int {
	s := ""
	for !p.slexer.isEOF() {
		c := p.slexer.peek()
		if ('a' <= c && c <= 'z') || 'A' <= c && c <= 'Z' || c == '_' || c == '.' || ('0' <= c && c <= '9') {
			s += string(c)
			p.slexer.next()
			continue
		}
		break
	}
	// Reserved Words
	lval.token = p.newToken(s)
	switch s {
	case "TIME", "Time":
		return TIME
	case "System.TimeSignature", "TimeSignature", "TimeSig":
		return TIME_SIG
	case "INT", "Int":
		return INT
	case "STR", "Str":
		return STR
	case "SUB", "Sub":
		return SUB
	}
	return WORD
}

func isFlag(c rune) bool { // 演算子か
	return c == '+' || c == '-' ||
		c == '*' || c == '/' || c == '%' ||
		c == '(' || c == ')' || c == '=' ||
		c == '#' || c == '"' || c == '\'' ||
		c == '>' || c == '<' || c == '[' || c == ']' ||
		c == '{' || c == '}' || c == ':' || c == '.' ||
		c == '!'
}

func isDigit(c rune) bool { // 数字か
	return '0' <= c && c <= '9'
}

func (p *Lexer) newToken(label string) Token {
	return Token{
		label:    label,
		lineinfo: p.slexer.getLineInfo(),
	}
}

// エラー報告用
func (p *Lexer) Error(e string) {
	tok := ""
	if p.lastToken != nil {
		tok = p.lastToken.label
	}
	e = strings.ReplaceAll(e, "syntax error", "文法エラー")
	e = strings.ReplaceAll(e, "unexpected $unk", "解析できない文字がありました。")
	msg := fmt.Sprintf("[ERROR] (%d) [%s] %s", CurLine.LineNo+1, tok, e)
	err := fmt.Errorf(msg)
	p.parseError = err
}
