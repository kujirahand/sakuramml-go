package parser

import (
	"fmt"
	"sakuramml/node"
	"sakuramml/token"
	"sakuramml/utils"
	"strconv"
)

// Parser struct
type Parser struct {
	desk token.Desk
	Top  *node.Node
	Last *node.Node
}

// Init parser
func (p *Parser) Init(tokens token.Tokens) {
	p.desk = token.NewDesk(tokens)
	nop := node.NewNop()
	p.Top = nop
	p.Last = p.Top
}

// NewParser func
func NewParser(tokens token.Tokens) *Parser {
	p := Parser{}
	p.Init(tokens)
	return &p
}

func (p *Parser) appendNode(n *node.Node) {
	if n == nil {
		return
	}
	p.Last.Next = n
	p.Last = n
}

func (p *Parser) readNoteOn(t *token.Token) (*node.Node, error) {
	ex := node.ExDataNoteOn{}
	n := node.NewNoteOn(t.Label, &ex)
	// sharp or flat
	for {
		if p.desk.IsLabel("+") || p.desk.IsLabel("#") || p.desk.IsLabel("♯") {
			ex.NoteShift++
			p.desk.Next()
			continue
		}
		if p.desk.IsLabel("-") || p.desk.IsLabel("♭") {
			ex.NoteShift--
			p.desk.Next()
			continue
		}
		break
	}
	// length ?
	if p.desk.IsType(token.Number) || p.desk.IsLabel("^") {
		nLen, err := p.readLength()
		if err != nil {
			return n, err
		}
		ex.Length = nLen
	}
	return n, nil
}

func (p *Parser) readRest(t *token.Token) (*node.Node, error) {
	ex := node.ExDataNoteOn{}
	n := node.NewNoteOn(t.Label, &ex)
	// length ?
	if p.desk.IsType(token.Number) || p.desk.IsLabel("^") {
		nLen, err := p.readLength()
		if err != nil {
			return n, err
		}
		ex.Length = nLen
	}
	return n, nil
}

func (p *Parser) readValue() (*node.Node, error) {
	if p.desk.IsLabel("=") {
		p.desk.Next()
	}
	ct := p.desk.Peek()
	if p.desk.IsType(token.Number) {
		nn := node.NewNumber(ct.Label)
		p.desk.Next()
		return nn, nil
	}
	return nil, fmt.Errorf("not implement : %s", ct.Label)
}

func (p *Parser) read1pCmd(t *token.Token, ntype node.NType) (*node.Node, error) {
	opt := ""
	for p.desk.IsLabel("%") || p.desk.IsLabel("+") || p.desk.IsLabel("-") {
		opt += p.desk.Next().Label
	}
	// read param
	if p.desk.IsLabel("(") { // skip ParenR
		p.desk.Next()
	}
	no, err := p.readValue()
	if err != nil {
		return nil, fmt.Errorf("%s : %s value invalid", t.Label, ntype)
	}
	if p.desk.IsLabel(")") { // skip ParenR
		p.desk.Next()
	}
	// process command
	switch ntype {
	case node.SetTrack:
		return node.NewSetTrack(no, opt), nil
	case node.SetOctave:
		return node.NewSetOctave(no, opt), nil
	case node.SetQgate:
		return node.NewSetQgate(no, opt), nil
	case node.SetVelocity:
		return node.NewSetVelocity(no, opt), nil
	case node.SetTempo:
		return node.NewSetTempo(no, opt), nil
	default:
		return nil, fmt.Errorf("System Error : No command : %s", ntype)
	}
}

func (p *Parser) readVoice(t *token.Token) (*node.Node, error) {
	// read param
	if p.desk.IsLabel("(") { // skip ParenR
		p.desk.Next()
	}
	no, err := p.readValue()
	if err != nil {
		return nil, fmt.Errorf("%s no invalid", t.Label)
	}
	// fix no
	msb := -1
	lsb := -1
	// msb
	if p.desk.IsLabel(",") {
		p.desk.Next()
		if !p.desk.IsType(token.Number) {
			return nil, fmt.Errorf("%s MSB no invalid should be number", t.Label)
		}
		msb, _ = strconv.Atoi(p.desk.Next().Label)
		msb = utils.MidiRange(msb)
		if p.desk.IsLabel(",") {
			p.desk.Next()
			if !p.desk.IsType(token.Number) {
				return nil, fmt.Errorf("%s MSB,LSB no invalid should be number", t.Label)
			}
			lsb, _ = strconv.Atoi(p.desk.Next().Label)
			lsb = utils.MidiRange(lsb)
		}
	}
	if p.desk.IsLabel(")") { // skip ParenR
		p.desk.Next()
	}
	nVoice := node.NewSetPC(no, msb, lsb)
	return nVoice, nil
}

func (p *Parser) readLength() (*node.Node, error) {
	nTop := node.NewNop()
	nLast := nTop
	loopc := 0
	for p.desk.HasNext() {
		// Number or Base(TrackLength)
		nNum := node.NewGetTrackLength()
		if p.desk.IsType(token.Number) {
			nValue, _ := p.readValue()
			nNum = node.NewNLenToStep(nValue)
		}
		res := nNum
		// Dot
		dotCount := 0
		dotRate := 1.0
		dotSum := 1.0
		for p.desk.IsLabel(".") {
			p.desk.Next()
			dotCount++
			dotRate = dotRate / 2.0
			dotSum += dotRate
		}
		if dotCount > 0 {
			nDot := node.NewLengthDot(nNum)
			nDot.ExData = float64(dotSum)
			res = nDot
		}
		nLast.Next = res
		nLast = nLast.Next
		// print("loop=", loopc, "\n", nLast.ToStringAll(), "\n")
		loopc++
		// Next
		if p.desk.IsLabel("^") {
			p.desk.Next()
			continue
		}
		break
	}
	if nTop == nLast {
		return node.NewGetTrackLength(), nil
	}
	// print("@@@\n")
	// fmt.Println(nTop.ToStringAll())
	nodeLength := node.NewLength()
	nodeLength.NValue = nTop
	return nodeLength, nil
}

func (p *Parser) readSetLength() (*node.Node, error) {
	if !p.desk.IsType(token.Number) {
		return nil, fmt.Errorf("l command need number")
	}
	nodeLength, err := p.readLength()
	if err != nil {
		return nil, err
	}
	return node.NewSetLength(nodeLength), nil
}

func (p *Parser) readWord() (*node.Node, error) {
	t := p.desk.Next()
	switch t.Label {
	case "c", "ド":
		return p.readNoteOn(t)
	case "d", "レ":
		return p.readNoteOn(t)
	case "e", "ミ":
		return p.readNoteOn(t)
	case "f", "フ":
		return p.readNoteOn(t)
	case "g", "ソ":
		return p.readNoteOn(t)
	case "a", "ラ":
		return p.readNoteOn(t)
	case "b", "シ":
		return p.readNoteOn(t)
	case "r", "ン", "ッ":
		t.Label = "r"
		return p.readRest(t)
	case "l":
		return p.readSetLength()
	case "o":
		return p.read1pCmd(t, node.SetOctave)
	case "v":
		return p.read1pCmd(t, node.SetVelocity)
	case "@", "Voice", "VOICE":
		return p.readVoice(t)
	case "TR", "Track":
		return p.read1pCmd(t, node.SetTrack)
	case "Tempo":
		return p.read1pCmd(t, node.SetTempo)
	}
	return nil, fmt.Errorf("Unknown Word : %s", t.Label)
}

func (p *Parser) readFlag() (*node.Node, error) {
	t := p.desk.Next()
	switch t.Label {
	case ">", "↑":
		return node.NewSetOctave(nil, "++"), nil
	case "<", "↓":
		return node.NewSetOctave(nil, "--"), nil
	case "@":
		return p.readVoice(t)
	}
	return nil, fmt.Errorf("Unknown Word : %s", t.Label)
}

// Parse func
func (p *Parser) Parse() (*node.Node, error) {
	var e error
	for p.desk.HasNext() {
		t := p.desk.Peek()
		if t.Type == token.Word {
			nn, err := p.readWord()
			if err != nil {
				return nil, err
			}
			p.appendNode(nn)
			continue
		}
		if t.Type == token.Flag {
			nn, err := p.readFlag()
			if err != nil {
				return nil, err
			}
			p.appendNode(nn)
			continue
		}
		e = fmt.Errorf("[ERROR] (%d) not implements : %s ", p.desk.Peek().Line, p.desk.Peek().Label)
		return p.Top, e
	}
	return p.Top, nil
}

// Parse convert to AST
func Parse(tokens token.Tokens) (*node.Node, error) {
	return NewParser(tokens).Parse()
}
