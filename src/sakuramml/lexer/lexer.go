package lexer

import (
    "log"
    "sakuramml/token"
)

type Lexer struct {
    input   []rune
    index   int
    tokens  []*token.Token
}

func (l *Lexer) Init(src string) {
    l.index = 0
    l.input = []rune(src)
    l.tokens = []*token.Token{}
}

func (l *Lexer) HasNext() bool {
    return (l.index < len(l.input))
}

func (l *Lexer) Split() {
    for l.HasNext() {
        // fmt.Println(l.index, ":", l.Peek())
        l.readOne()
    }
}

func (l *Lexer) Peek() rune {
    if l.index >= len(l.input) {
        return rune(0)
    }
    return l.input[l.index]
}

func (l *Lexer) Next() rune {
    var ch = rune(0)
    if l.index < len(l.input) {
        ch = l.input[l.index]
    }
    l.index += 1
    return ch
}

func isSpace(c rune) bool {
    return c == rune(' ')  || c == rune('\t') || c == rune('\r') || c == rune('\n')
}

func (l *Lexer) SkipSpace() {
    for {
        ch := l.Peek()
        if isSpace(ch) {
            l.index++
            continue
        }
        break
    }
}

func IsLower(c rune) bool {
    return rune('a') <= c && c <= rune('z')
}
func IsUpper(c rune) bool {
    return rune('A') <= c && c <= rune('Z')
}
func IsDigit(c rune) bool {
    return rune('0') <= c && c <= rune('9')
}
func (l *Lexer) readWord() string {
    var s = ""
    // Upper Char
    fc := l.Next()
    if !IsUpper(fc) { return string(fc) }
    for l.HasNext() {
        ch := l.Peek()
        if IsUpper(ch) || IsLower(ch) || IsDigit(ch) || ch == rune('_') {
            s += string(ch)
        } else {
            break
        }
    }
    return s
}

func (l *Lexer) readOne() {
    l.SkipSpace()
    ch := l.Peek()
    if ch == rune(0) { return }
    if IsLower(ch) {
        t := token.Token{Type:token.WORD, Label:string(ch)}
        l.tokens = append(l.tokens, &t)
        l.Next()
        return
    }
    if IsUpper(ch) {
        w := l.readWord()
        t := token.Token{Type:token.WORD, Label:w}
        l.tokens = append(l.tokens, &t)
        return
    }
    log.Fatal("[ERROR] Unknown word: " + string(ch))
    l.Next()
}

func Lex(src string) []*token.Token {
    l := Lexer{}
    l.Init(src)
    l.Split()
    return l.tokens
}


