package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/kujirahand/sakuramml-go"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		ShowHelp()
		return
	}
	// Check command line Options
	opt := sakuramml.CompilerOptions{}
	for i, arg := range args {
		if i == 0 {
			continue
		} // exefile
		if arg == "--help" || arg == "-h" || arg == "?" {
			ShowHelp()
			return
		}
		if arg == "--debug" || arg == "-d" {
			opt.Debug = true
			continue
		}
		if arg == "--e" || arg == "-e" {
			opt.EvalMode = true
			continue
		}
		// EvalMode source
		if opt.EvalMode {
			if "" == opt.Source {
				opt.Source = arg
				continue
			}
		} else {
			// Check in out filename
			if "" == opt.Infile {
				opt.Infile = arg
				continue
			}
		}
		if "" == opt.Outfile {
			opt.Outfile = arg
			continue
		}
	}
	// No infile ?
	if opt.Source == "" && opt.Infile == "" {
		ShowHelp()
		return
	}
	// No outfile
	if opt.Outfile == "" {
		re := regexp.MustCompile("\\.mml$")
		mml := opt.Infile
		if mml == "" {
			mml = "a.mml"
		}
		out := re.ReplaceAllString(mml, ".mid")
		if mml == out {
			out += ".mid"
		}
		opt.Outfile = out
	}
	if opt.Debug {
		fmt.Println("Command line:", opt)
	}
	// load file
	if !opt.EvalMode {
		text, err := ioutil.ReadFile(opt.Infile)
		if err != nil {
			log.Fatal("[ERROR] Fail to load infile: " + opt.Infile)
		}
		opt.Source = string(text)
	}
	// run
	song, err := sakuramml.Compile(&opt)
	if err != nil {
		log.Fatal(err)
	}
	// save to file
	sakuramml.MidiSaveToFile(song, opt.Outfile)
	if opt.Debug {
		fmt.Printf("SaveToFile=%s\n", opt.Outfile)
	}
	fmt.Println("ok.")
}

// ShowHeader func
func ShowHeader() {
	fmt.Println("♪ sakuramml-go " + sakuramml.VERSION)
}

// ShowHelp func
func ShowHelp() {
	ShowHeader()
	fmt.Println("USAGE:")
	fmt.Println("  csakura (mmlfile) [--out=(midifile)]")
	fmt.Println("OPTIONS:")
	fmt.Println("  -h, --help     Show Help")
	fmt.Println("  -d, --debug    Debug mode")
	fmt.Println("  -e (mml)       eval mode")
}
