package sakuramml

import (
	"fmt"
)

const (
	// VERSION : sakuramml version
	VERSION = "0.0.1"
)

// CompilerOptions : Compiler CompilerOptions
type CompilerOptions struct {
	Debug    bool
	EvalMode bool
	Infile   string
	Source   string
	Outfile  string
}

// Eval func
func Eval(song *Song, src string) error {
	topNode, err := Parse(src, 0)
	if err != nil {
		return err
	}
	// run
	return CompilerRun(topNode, song)
}

// Run func
func CompilerRun(topNode *Node, song *Song) error {
	return nil
}

// Compile MML
func Compile(opt *CompilerOptions) (*Song, error) {
	// init
	songObj := NewSong()
	songObj.Debug = opt.Debug
	songObj.Eval = Eval // Set Eval Func
	// sutoton
	if opt.Debug {
		fmt.Println("--- sutoton ---")
	}
	src, err := SutotonConvert(opt.Source)
	if err != nil {
		return nil, err
	}
	// parse
	if opt.Debug {
		fmt.Println("--- parse ---")
	}
	node, err := Parse(src, 0)
	if err != nil {
		return nil, err
	}
	// exec
	if opt.Debug {
		fmt.Println("--- exec ---")
	}
	Run(node, songObj)
	return songObj, nil
}
