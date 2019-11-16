package lexer

import (
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
    for !l.HasNext() {
        l.Next()
    }
}

func (l *Lexer) Peek() rune {
    return l.input[l.index]
}
func (l *Lexer) Get() rune {
    ch := l.input[l.index]
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

func (l *Lexer) GetWord() {
}

func (l *Lexer) Next() {
    l.SkipSpace() 
}

func Lex(src string) []*token.Token {
    l := Lexer{}
    l.Init(src)
    l.Split()
    return l.tokens
}
