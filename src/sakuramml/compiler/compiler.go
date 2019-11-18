package compiler

import (
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
	song := song.Song{}
	song.Init()
	// lex
	tokens, err := lexer.Lex(opt.Source)
	if err != nil {
		return nil, err
	}
	parser.Parse(tokens)
	print(song.ToString())
	print(token.TokensToString(tokens))
	return &song, nil
}
