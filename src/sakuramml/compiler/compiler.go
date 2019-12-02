package compiler

import (
	"fmt"
	"sakuramml/lexer"
	"sakuramml/node"
	"sakuramml/parser"
	"sakuramml/song"
	"sakuramml/sutoton"
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


// Eval func
func Eval(song *song.Song, src string) error {
	// lex
	tokens, err := lexer.Lex(src)
	if err != nil {
		return err
	}
	// parse
	topNode, err := parser.Parse(tokens)
	if err != nil {
		return err
	}
	// run
	return Run(topNode, song)
}

// Run func
func Run(topNode *node.Node, song *song.Song) error {
	curNode := topNode
	for curNode != nil {
		// Run Node
		err := curNode.Exec(curNode, song)
		if err != nil {
			return err
		}
		// Force Change Node?
		if song.MoveNode != nil {
			curNode = song.MoveNode.(*node.Node)
			song.MoveNode = nil
			continue
		}
		curNode = curNode.Next
	}
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
	// lex
	if opt.Debug {
		fmt.Println("--- lex ---")
	}
	tokens, err := lexer.Lex(src)
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
	// fmt.Println(s.ToString())
	for {
		err := Run(topNode, songObj)
		if err != nil {
			return nil, err
		}
		break
	}
	return songObj, nil
}
