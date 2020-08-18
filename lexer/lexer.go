package lexer

import (
	"log"
	"github.com/kujirahand/sakuramml-go/token"
	"github.com/kujirahand/sakuramml-go/utils"
)

// Lexer struct
type Lexer struct {
	input  []rune
	index  int
	length int
	line   int
	tokens token.Tokens
}

// NewLexer func
func NewLexer(src string) *Lexer {
	l := Lexer{}
	l.SetInput(src)
	return &l
}

// SetInput func
func (l *Lexer) SetInput(src string) {
	l.index = 0
	l.line = 0
	l.input = []rune(src)
	l.length = len(l.input)
	l.tokens = token.Tokens{}
}

// HasNext : Check Next rune
func (l *Lexer) HasNext() bool {
	return (l.index < l.length)
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
	if l.index >= l.length {
		return rune(0)
	}
	return l.input[l.index]
}

// IsLabel func
func (l *Lexer) IsLabel(s string) bool {
	return utils.StrCompareKey(l.input, l.index, s)
}

// Next : Get current rune and inc index
func (l *Lexer) Next() rune {
	var ch = rune(0)
	if l.index < l.length {
		ch = l.input[l.index]
	}
	l.index++
	return ch
}

// Move : Move cursor
func (l *Lexer) Move(n int) rune {
	l.index += n
	if l.index < 0 {
		l.index = 0
	}
	return l.Peek()
}

// IsSpace is check white space rune
func IsSpace(c rune) bool {
	return c == rune(' ') || c == rune('\t') || c == rune('\r') || c == rune('\n')
}

// SkipSpace : skip space
func (l *Lexer) SkipSpace() {
	for {
		ch := l.Peek()
		if IsSpace(ch) {
			l.index++
			if ch == '\n' {
				l.line++
			}
			continue
		}
		break
	}
}

// readLineComment func
func (l *Lexer) readLineComment() string {
	if !l.IsLabel("//") {
		return ""
	}
	if l.IsLabel("///") {
		l.Move(3)
	}
	comment := ""
	for l.HasNext() {
		c := l.Next()
		if c == rune('\n') {
			l.line++
			break
		}
		comment += string(c)
	}
	return "/*" + comment + "*/"
}

// readRangeComment func ... could nest
func (l *Lexer) readRangeComment() string {
	if !l.IsLabel("/*") {
		return ""
	}
	comment := ""
	l.Move(2) // skip "/*"
	level := 1
	for l.HasNext() {
		if l.IsLabel("/*") {
			l.Move(2)
			comment += "/*"
			level++
			continue
		}
		if l.IsLabel("*/") {
			l.Move(2) // skip "*/"
			comment += "*/"
			level--
			if level == 0 {
				break
			}
		}
		c := l.Next()
		if c == '\n' {
			l.line++
		}
		comment += string(c)
	}
	return comment
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

func (l *Lexer) readMacro() string {
	if l.Peek() != '#' {
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
	// fmt.Printf("ch=%s, %d\n", string(ch), int(ch))
	// line comment ?
	if ch == rune('/') {
		// embed line comment
		if l.IsLabel("///") {
			l.appendToken(token.Comment, l.readLineComment())
			return
		}
		if l.IsLabel("//") {
			l.readLineComment() // Only Read, not append
			l.SkipSpace()
			l.readOne()
			return
		}
		// range comment
		if l.IsLabel("/*") {
			l.readRangeComment()
			l.SkipSpace()
			l.readOne()
			return
		}
	}
	// Multi Byte Rune
	if int(ch) > 0xFF { // ストトンで変換できなかったもの
		l.appendToken(token.Word, string(ch))
		l.Next()
		return
	}
	// Macro
	if ch == '#' {
		w := l.readMacro()
		l.appendToken(token.Macro, w)
		return
	}
	// Str
	if ch == '{' {
		str := l.readNestString()
		l.appendToken(token.String, str)
		return
	}
	// Lower Rune
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
	if IsFlag(ch) {
		l.appendToken(token.Flag, string(ch))
		l.Next()
		return
	}
	log.Fatal("[ERROR] Unknown char: " + string(ch))
	l.Next()
}

func (l *Lexer) readNestString() string {
	if l.Peek() != '{' {
		return ""
	}

	l.index++ // skip {
	level := 1
	res := ""
	for l.HasNext() {
		c := l.Next()
		if c == '}' {
			level--
			if level <= 0 {
				break
			}
		} else if c == '{' {
			level++
		} else if c == '\n' {
			l.line++
		}
		res += string(c)
	}
	return res
}

func (l *Lexer) appendToken(tt token.TType, label string) {
	t := token.Token{Type: tt, Label: label, Line: l.line}
	l.tokens = append(l.tokens, t)
}

// Lex : split to tokens
func Lex(src string) (token.Tokens, error) {
	return NewLexer(src).Split()
}
