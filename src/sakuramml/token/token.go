package token

import (
	"fmt"
)

const (
	// Word : Token type
	Word = "word"
	// Number : Token type
	Number = "number"
	// Flag : Token type
	Flag = "flag"
	// ParenL : Token type
	ParenL = "("
	// ParenR : Token type
	ParenR = ")"
	// BracketL : Token type
	BracketL = "["
	// BracketR : Token type
	BracketR = "]"
)

// Token struct
type Token struct {
	Type  string
	Label string
}

// Tokens Slice
type Tokens []Token

// TokensToString for Debug
func TokensToString(tokens Tokens) string {
	s := ""
	for i, t := range tokens {
		s += fmt.Sprintf("%3d: %5s %s\n", i, t.Type, t.Label)
	}
	return s
}
