package token

import (
    "fmt"
)

type TokenType string

const (
    WORD = "word"
)

type Token struct {
    Type    TokenType
    Label   string
}

type Tokens []*Token

func TokensToString(tokens Tokens) string {
    s := ""
    for i, t := range tokens {
        s += fmt.Sprintf("%3d: %5s %s\n", i, t.Type, t.Label)
    }
    return s
}


