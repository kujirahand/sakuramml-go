package main

import (
    "os"
    "fmt"
    "regexp"
    "log"
    "io/ioutil"
    "sakuramml/compiler"
    "sakuramml/midi"
)

func main() {
    args := os.Args
    if len(args) <= 1 {
        ShowHelp()
        return
    }
    // Check command line Options
    opt := compiler.Options{}
    for i, arg := range args {
        if i == 0 { continue } // exefile
        if arg == "--help" || arg == "-h" || arg == "?" {
            ShowHelp()
            return
        }
        if arg == "--debug" || arg == "-d" {
            opt.Debug = true
            continue
        }
        // Check in out filename
        if "" == opt.Infile {
            opt.Infile = arg
            continue
        }
        if "" == opt.Outfile {
            opt.Outfile = arg
            continue
        }
    }
    // No infile ?
    if opt.Infile == "" {
        ShowHelp()
        return
    }
    // No outfile
    if opt.Outfile == "" {
        re := regexp.MustCompile("\\.mml$")
        mml := opt.Infile
        out := re.ReplaceAllString(mml, ".mid")
        if mml == out { out += ".mid" }
        opt.Outfile = out
    }
    if opt.Debug {
        fmt.Println("Command line:", opt)
    }
    // load file
    text, err := ioutil.ReadFile(opt.Infile)
    if err != nil {
        log.Fatal("[ERROR] Fail to load infile: " + opt.Infile)
    }
    opt.Source = string(text)
    // run
    song, err := compiler.Compile(&opt)
    if err != nil {
        log.Fatal(err)
    }
    f, err := os.Create(opt.Outfile)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    midi.Save(song, f)
}

func ShowHeader() {
    fmt.Println("♪ sakuramml-go " + compiler.VERSION)
}

func ShowHelp() {
    ShowHeader()
    fmt.Println("USAGE:")
    fmt.Println("  csakura (mmlfile) [--out=(midifile)]")
    fmt.Println("OPTIONS:")
    fmt.Println("  -h, --help     Show Help")
    fmt.Println("  -d, --debug    Debug mode")
}

