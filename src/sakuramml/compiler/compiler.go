package compiler

import (
    "sakuramml/song"
    "sakuramml/lexer"
    "sakuramml/token"
)

const (
    VERSION = "0.0.1"
)

type Options struct {
    Debug bool
    Infile string
    Source string
    Outfile string
}

func Compile(opt *Options) bool {
    song := song.Song{}
    song.Init()
    tokens := lexer.Lex(opt.Source)
    print(song.ToString())
    print(token.TokensToString(tokens))
    return true
}






