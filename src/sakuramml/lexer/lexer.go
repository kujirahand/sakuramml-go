package lexer

import (
	"log"
	"sakuramml/token"
)

// Lexer struct
type Lexer struct {
	input  []rune
	index  int
	tokens token.Tokens
}

// Init : Initialize Lexer struct
func (l *Lexer) Init(src string) {
	l.index = 0
	l.input = []rune(src)
	l.tokens = token.Tokens{}
}

// NewLexer func
func NewLexer(src string) *Lexer {
	l := Lexer{}
	l.Init(src)
	return &l
}

// HasNext : Check Next rune
func (l *Lexer) HasNext() bool {
	return (l.index < len(l.input))
}

// Split : Get tokens
func (l *Lexer) Split() (token.Tokens, error) {
	for l.HasNext() {
		l.readOne()
	}
	return l.tokens, nil
}

// Peek current rune
func (l *Lexer) Peek() rune {
	if l.index >= len(l.input) {
		return rune(0)
	}
	return l.input[l.index]
}

// Next : Get current rune and inc index
func (l *Lexer) Next() rune {
	var ch = rune(0)
	if l.index < len(l.input) {
		ch = l.input[l.index]
	}
	l.index++
	return ch
}

// IsSpace is check whilte space rune
func IsSpace(c rune) bool {
	return c == rune(' ') || c == rune('\t') || c == rune('\r') || c == rune('\n')
}

// SkipSpace : skip space
func (l *Lexer) SkipSpace() {
	for {
		ch := l.Peek()
		if IsSpace(ch) {
			l.index++
			continue
		}
		break
	}
}

// IsLower : Is rune lower case?
func IsLower(c rune) bool {
	return rune('a') <= c && c <= rune('z')
}

// IsUpper : Is rune upper case?
func IsUpper(c rune) bool {
	return rune('A') <= c && c <= rune('Z')
}

// IsDigit : Is rune Digit?
func IsDigit(c rune) bool {
	return rune('0') <= c && c <= rune('9')
}

// IsFlag : Is rune Flag?
func IsFlag(c rune) bool {
	return rune(0x21) <= c && c <= rune(0x2F) ||
		rune(0x3A) <= c && c <= rune(0x40) ||
		rune(0x5B) <= c && c <= rune(0x60) ||
		rune(0x7B) <= c && c <= rune(0x7E)
}

func (l *Lexer) readWord() string {
	if !IsUpper(l.Peek()) {
		return ""
	}
	var s = ""
	for l.HasNext() {
		ch := l.Peek()
		if IsUpper(ch) || IsLower(ch) || IsDigit(ch) || ch == rune('_') {
			s += string(ch)
			l.Next()
		} else {
			break
		}
	}
	return s
}

func (l *Lexer) readNumber() string {
	if !IsDigit(l.Peek()) {
		return ""
	}
	var s = ""
	if l.Peek() == rune('$') {
		s += "0x"
	}
	for l.HasNext() {
		ch := l.Peek()
		if IsDigit(ch) {
			s += string(ch)
			l.Next()
		} else {
			break
		}
	}
	return s
}

func (l *Lexer) readOne() {
	l.SkipSpace()
	ch := l.Peek()
	if ch == rune(0) {
		return
	}
	if IsLower(ch) {
		l.appendToken(token.Word, string(ch))
		l.Next()
		return
	}
	if IsUpper(ch) {
		w := l.readWord()
		l.appendToken(token.Word, w)
		return
	}
	if IsDigit(ch) || ch == rune('$') {
		num := l.readNumber()
		l.appendToken(token.Number, num)
		return
	}
	switch ch {
	case rune('@'): // Voice
		l.appendToken(token.Word, string(ch))
		l.Next()
		return
	case rune('('):
		l.appendToken(token.ParenL, string(ch))
		l.Next()
		return
	case rune(')'):
		l.appendToken(token.ParenR, string(ch))
		l.Next()
		return
	case rune('['):
		l.appendToken(token.BracketL, string(ch))
		l.Next()
		return
	case rune(']'):
		l.appendToken(token.BracketR, string(ch))
		l.Next()
		return
	default:
		if IsFlag(ch) {
			l.appendToken(token.Flag, string(ch))
			l.Next()
			return
		}
		//
	}
	log.Fatal("[ERROR] Unknown word: " + string(ch))
	l.Next()
}
func (l *Lexer) appendToken(tt token.TokenType, label string) {
	t := token.Token{Type: tt, Label: label}
	l.tokens = append(l.tokens, t)
}

// Lex : split to tokens
func Lex(src string) (token.Tokens, error) {
	return NewLexer(src).Split()
}
