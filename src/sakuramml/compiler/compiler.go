package compiler

import (
    "sakuramml/song"
    "sakuramml/lexer"
    "sakuramml/token"
    "sakuramml/parser"
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






