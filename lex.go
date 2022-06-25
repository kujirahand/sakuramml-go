package sakuramml

import (
	"fmt"
	"strconv"
)

// Lexer 最低限必要な構造体を定義
type Lexer struct {
	src    []rune
	index  int
	lineno int
	fileno int
	result *Node
}

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

// Parse : コードを実行
func Parse(code string) (*Node, error) {
	lexer := &Lexer{
		src:   []rune(code + "\n"),
		index: 0,
	}
	yyDebug = 1
	yyErrorVerbose = true
	yyParse(lexer)
	return lexer.result, nil
}

// Lex : トークンを一つずつ返す
func (p *Lexer) Lex(lval *yySymType) int {
	tok := p.getToken(lval)
	fmt.Printf("[%s]", getTokenName(tok))
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
	p.skipSpace()
	// Set LineInfo
	CurLine.LineNo = p.lineno
	CurLine.FileNo = p.fileno
	// Check EOF
	if p.isEOF() {
		return 0 // end
	}
	c := p.peek()
	// LF ?
	if c == '\n' {
		lval.token = p.newToken("\n")
		p.lineno++
		p.next()
		return LF
	}
	// NUMBER ?
	if isDigit(c) {
		return p.lexNumber(lval)
	}
	// 演算子 ?
	if isFlag(c) {
		p.next()
		lval.token = p.newToken(string(c))
		return int(c)
	}
	// 小文字コマンド
	if 'a' <= c && c <= 'z' {
		lval.token = p.newToken(string(c))
		p.next()
		return int(c)
	}
	// 大文字コマンド
	if 'A' <= c && c <= 'Z' || c == '_' {
		return p.lexWord(lval)
	}
	return -1
}

func (p *Lexer) lexNumber(lval *yySymType) int {
	s := ""
	for !p.isEOF() {
		c := p.peek()
		if isDigit(c) || c == '.' {
			s += string(c)
			p.next()
			continue
		}
		break
	}
	lval.token = p.newToken(s)
	lval.token.value, _ = strconv.ParseFloat(s, 64)
	return NUMBER
}

func (p *Lexer) lexWord(lval *yySymType) int {
	s := ""
	for !p.isEOF() {
		c := p.peek()
		if ('a' <= c && c <= 'z') || 'A' <= c && c <= 'Z' || c == '_' {
			s += string(c)
			p.next()
			continue
		}
		break
	}
	lval.token = p.newToken(s)
	return WORD
}

func (p *Lexer) isEOF() bool {
	for p.index >= len(p.src) {
		return true
	}
	return false
}

func (p *Lexer) peek() rune {
	return p.src[p.index]
}

func (p *Lexer) next() {
	p.index++
}

func (p *Lexer) skipSpace() {
	for !p.isEOF() {
		c := p.peek()
		if c == ' ' || c == '\t' || c == '\r' {
			p.next()
			continue
		}
		break
	}
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

// エラー報告用
func (p *Lexer) Error(e string) {
	fmt.Println("[error] " + e)
}

func (p *Lexer) getLineInfo() LineInfo {
	return LineInfo{LineNo: p.lineno, FileNo: p.fileno}
}
func (p *Lexer) newToken(label string) Token {
	return Token{
		label:    label,
		lineinfo: p.getLineInfo(),
	}
}
