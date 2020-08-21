package compiler

import (
	"fmt"

	"github.com/kujirahand/sakuramml-go/mmlparser"
	"github.com/kujirahand/sakuramml-go/song"
	"github.com/kujirahand/sakuramml-go/sutoton"
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

// Eval func
func Eval(song *song.Song, src string) error {
	topNode, err := mmlparser.Parse(src)
	if err != nil {
		return err
	}
	// run
	return Run(topNode, song)
}

// Run func
func Run(topNode *mmlparser.Node, song *song.Song) error {
	return nil
}

// Compile MML
func Compile(opt *Options) (*song.Song, error) {
	// init
	songObj := song.NewSong()
	songObj.Debug = opt.Debug
	songObj.Eval = Eval // Set Eval Func
	// sutoton
	if opt.Debug {
		fmt.Println("--- sutoton ---")
	}
	src, err := sutoton.Convert(opt.Source)
	if err != nil {
		return nil, err
	}
	// parse
	if opt.Debug {
		fmt.Println("--- parse ---")
	}
	node, err := mmlparser.Parse(src)
	if err != nil {
		return nil, err
	}
	// exec
	if opt.Debug {
		fmt.Println("--- exec ---")
	}
	mmlparser.Run(node, songObj)
	return songObj, nil
}
