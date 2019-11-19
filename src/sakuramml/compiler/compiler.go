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
	Debug    bool
	EvalMode bool
	Infile   string
	Source   string
	Outfile  string
}

// Compile MML
func Compile(opt *Options) (*song.Song, error) {
	// init
	s := song.NewSong()
	s.Debug = opt.Debug
	// lex
	if opt.Debug {
		fmt.Println("--- lex ---")
	}
	tokens, err := lexer.Lex(opt.Source)
	if err != nil {
		return nil, err
	}
	if opt.Debug {
		fmt.Println(token.TokensToString(tokens, " "))
	}
	// parse
	if opt.Debug {
		fmt.Println("--- parse ---")
	}
	topNode, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}
	if opt.Debug {
		fmt.Println(topNode.ToStringAll())
	}
	// exec
	if opt.Debug {
		fmt.Println("--- exec ---")
	}
	curNode := topNode
	for curNode != nil {
		// if opt.Debug { fmt.Println(curNode.Type) }
		curNode.Exec(curNode, s)
		curNode = curNode.Next
	}
	// fmt.Println(s.ToString())
	return s, nil
}
