package sakuramml

import (
	"fmt"
	"regexp"
)

const (
	// VERSION : sakuramml version
	VERSION = "0.0.1"
)

var SakuraDebug bool = true

func SakuraLog(msg string) {
	if SakuraDebug {
		re := regexp.MustCompile(`\s+$`)
		msg = re.ReplaceAllString(msg, "")
		fmt.Println("[LOG] " + msg)
	}
}

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
	return SakuraRun(topNode, song)
}

// Compile MML
func Compile(opt *CompilerOptions) (*Song, error) {
	// init
	songObj := NewSong()
	songObj.Debug = opt.Debug
	SakuraDebug = songObj.Debug
	songObj.Eval = Eval // Set Eval Func
	// sutoton
	SakuraLog("--- Sutoton ---")
	src, err := SutotonConvert(opt.Source)
	// SakuraLog(src)
	if err != nil {
		return nil, err
	}
	// parse
	SakuraLog("--- Parse ---")
	node, err := Parse(src, 0)
	if err != nil {
		return nil, err
	}
	// exec
	SakuraLog("--- Run ---")
	SakuraRun(node, songObj)
	return songObj, nil
}
