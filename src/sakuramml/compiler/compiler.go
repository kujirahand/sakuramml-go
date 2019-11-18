package compiler

import (
	"fmt"
	"sakuramml/lexer"
	"sakuramml/parser"
	"sakuramml/song"
	"sakuramml/token"
)

const (
	// VERSION : sakuramml version
	VERSION = "0.0.1"
)

// Options : Compiler Options
type Options struct {
	Debug   bool
	Infile  string
	Source  string
	Outfile string
}

// Compile MML
func Compile(opt *Options) (*song.Song, error) {
	// init
	s := song.NewSong()
	// lex
	tokens, err := lexer.Lex(opt.Source)
	if err != nil {
		return nil, err
	}
	fmt.Println(token.TokensToString(tokens, " "))
	// parse
	topNode, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}
	// exec
	curNode := topNode
	for curNode != nil {
		curNode.Exec(curNode, s)
		curNode = curNode.Next
	}
	fmt.Println(s.ToString())
	return s, nil
}
