package compiler

import (
    "sakuramml/song"
    "sakuramml/lexer"
)

const (
    VERSION = "0.0.1"
)

type Options struct {
    Debug bool
    Infile string
    Outfile string
}

func Compile(opt *Options) bool {
    song := song.Song{}
    song.Init()
    print(song.ToString())
    return true
}






