package parser

import (
	"fmt"
	"sakuramml/node"
	"sakuramml/token"
	"sakuramml/track"
	"sakuramml/utils"
	"sakuramml/variable"
	"strconv"
)

// Parser struct
type Parser struct {
	desk token.Desk
	harmonyStack []*node.Node
	Top  *node.Node
	Last *node.Node
	Variable *variable.Variable // temporary variable def
}

// NewParser func
func NewParser(tokens token.Tokens) *Parser {
	p := Parser{}
	p.desk = token.NewDesk(tokens)
	p.harmonyStack = []*node.Node{}
	nop := node.NewNop()
	p.Top = nop
	p.Last = p.Top
	p.Variable = variable.NewVariable()
	return &p
}

func (p *Parser) readWord() (*node.Node, error) {
	t := p.desk.Next()
	switch t.Label {
	case "c":
		return p.readNoteOn(t)
	case "d":
		return p.readNoteOn(t)
	case "e":
		return p.readNoteOn(t)
	case "f":
		return p.readNoteOn(t)
	case "g":
		return p.readNoteOn(t)
	case "a":
		return p.readNoteOn(t)
	case "b":
		return p.readNoteOn(t)
	case "n":
		return p.readNoteOn(t)
	case "r":
		return p.readRest(t)
	case "l":
		return p.readSetLength()
	case "o":
		return p.read1pCmd(t, node.SetOctave)
	case "q":
		return p.read1pCmd(t, node.SetQgate)
	case "v":
		return p.read1pCmd(t, node.SetVelocity)
	case "p":
		return p.read1pCmd(t, node.SetPitchBend)
	case "@", "Voice", "VOICE":
		return p.readVoice(t)
	case "TR", "Track", "TRACK":
		return p.read1pCmd(t, node.SetTrack)
	case "Tempo", "TEMPO", "BPM":
		return p.read1pCmd(t, node.SetTempo)
	case ">", "↑":
		return node.NewSetOctave(nil, "++"), nil
	case "<", "↓":
		return node.NewSetOctave(nil, "--"), nil
	case "[":
		return p.readLoopBegin()
	case "]":
		return p.readLoopEnd()
	case ":":
		return p.readLoopBreak()
	case "|":
		return node.NewNop(), nil
	case "Int", "INT":
		return p.readInt()
	case "Str", "STR":
		return p.readStr()
	case "Print", "PRINT":
		return p.readPrint()
	default:
		// eval
		varName := t.Label
		if p.Variable.Exists(varName) {
			return node.NewStrEval(varName), nil
		}
	}
	return nil, fmt.Errorf("[Error] (%d) Unknown Word : %s", t.Line + 1, t.Label)
}
func (p *Parser) readPrint() (*node.Node, error) {
	pnode, err := p.readValue()
	if err != nil {
		return nil, err
	}
	n := node.NewPrint(pnode)
	return n, nil
}

// Parse func
func (p *Parser) Parse() (*node.Node, error) {
	var nod *node.Node
	var err error
	for p.desk.HasNext() {
		tok := p.desk.Peek()
		// fmt.Printf("parse %v\n", *tok)
		switch tok.Type {
		case token.Word, token.Flag:
			nod, err = p.readWord()
			if err != nil {
				return nil, err
			}
			p.appendNode(nod)
			continue
		case token.Macro:
			nod, err = p.readMacro()
			if err != nil {
				return nil, err
			}
			p.appendNode(nod)
			continue
		case token.Comment:
			p.appendNode(node.NewComment(tok.Label))
			p.desk.Next()
			continue
		}
		err = fmt.Errorf("[ERROR] (%d) Parser not implemented : %s ",
			p.desk.Peek().Line, p.desk.Peek().Label)
		return p.Top, err
	}
	return p.Top, nil
}

func (p *Parser) appendNode(n *node.Node) {
	if n == nil {
		return
	}
	p.Last.Next = n
	p.Last = n
}

// readNoteOn func
func (p *Parser) readNoteOn(t *token.Token) (*node.Node, error) {
	var err error
	ex := node.ExDataNoteOn{}
	n := node.NewNoteOn(t.Label, &ex)
	isHarmony := false
	// n command ?
	if t.Label == "n" {
		noteNoNode, err := p.readValue()
		if err != nil {
			return nil, fmt.Errorf("[ERRPR] (%d) Failed to read n command NoteNo", t.Line)
		}
		ex.NoteNo = noteNoNode
		if p.desk.IsLabel(",") {
			p.desk.Next()
		}
	}
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
		if p.desk.IsLabel("0") {
			p.desk.Next() // skip "0"
			p.harmonyStack = append(p.harmonyStack, n)
			isHarmony =true
		} else {
			nLen, err := p.readLength()
			if err != nil {
				return n, err
			}
			ex.Length = nLen
		}
	}
	// qgate ?
	if p.desk.IsLabel(",") {
		p.desk.Next()
		if !p.desk.IsLabel(",") { // 省略がなければ暫定qを読む
			ex.QgateMode = track.QgateModeRate
			if p.desk.IsLabel("+") || p.desk.IsLabel("-") {
				qf := p.desk.Peek().Label
				ex.QgateOpt = rune(qf[0])
				p.desk.Next()
			}
			if p.desk.IsLabel("%") {
				ex.QgateMode = track.QgateModeStep
				p.desk.Next()
			}
			ex.Qgate, err = p.readValue()
			if err != nil {
				return nil, fmt.Errorf("NoteOn(l,[q]) Error")
			}
		}
		// velocity ?
		if p.desk.IsLabel(",") {
			p.desk.Next()
			if !p.desk.IsLabel(",") { // 省略がなければvを読む
				if p.desk.IsLabel("+") || p.desk.IsLabel("-") {
					ex.VelocityOpt = rune(p.desk.Next().Label[0])
				}
				ex.Velocity, err = p.readValue()
				if err != nil {
					return nil, fmt.Errorf("NoteOn(l, q, [v]) Error")
				}
			}
		}
	}
	if !isHarmony && len(p.harmonyStack) > 0{

	}
	return n, nil
}

func (p *Parser) readRest(t *token.Token) (*node.Node, error) {
	ex := node.ExDataNoteOn{}
	n := node.NewNoteOn("r", &ex)
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

func (p *Parser) readLoopBegin() (*node.Node, error) {
	var loopValue *node.Node
	var err error
	// loop counter ?
	if p.desk.IsType(token.Number) || p.desk.IsLabel("(") {
		loopValue, err = p.readValue()
		if err != nil {
			return nil, err
		}
	} else {
		// no counter => 2
		loopValue = node.NewNumberInt(2)
	}

	nodeLoopBegin := node.NewLoopBegin(loopValue)
	return nodeLoopBegin, nil
}
func (p *Parser) readLoopEnd() (*node.Node, error) {
	return node.NewLoopEnd(), nil
}
func (p *Parser) readLoopBreak() (*node.Node, error) {
	return node.NewLoopBreak(), nil
}

func (p *Parser) calcExpr() (*node.Node, error) {
	left, err := p.calcTerm()
	if err != nil {
		return nil, err
	}
	for p.desk.HasNext() {
		if p.desk.IsLabel("+") {
			p.desk.Next()
			right, err := p.calcTerm()
			if err != nil {
				return nil, err
			}
			left = node.NewCalcAdd(left, right)
			continue
		}
		if p.desk.IsLabel("-") {
			p.desk.Next()
			right, err := p.calcTerm()
			if err != nil {
				return nil, err
			}
			left = node.NewCalcSub(left, right)
			continue
		}
		break
	}
	return left, nil
}

func (p *Parser) calcTerm() (*node.Node, error) {
	left, err := p.calcFactor()
	if err != nil {
		return nil, err
	}
	for p.desk.HasNext() {
		if p.desk.IsLabel("*") {
			p.desk.Next()
			right, err := p.calcFactor()
			if err != nil {
				return nil, err
			}
			left = node.NewCalcMul(left, right)
			continue
		}
		if p.desk.IsLabel("/") {
			p.desk.Next()
			right, err := p.calcFactor()
			if err != nil {
				return nil, err
			}
			left = node.NewCalcDiv(left, right)
			continue
		}
		if p.desk.IsLabel("%") {
			p.desk.Next()
			right, err := p.calcFactor()
			if err != nil {
				return nil, err
			}
			left = node.NewCalcMod(left, right)
			continue
		}
		break
	}
	return left, nil
}

func (p *Parser) calcFactor() (*node.Node, error) {
	t := p.desk.Peek()
	// ( .. )
	if p.desk.IsLabel("(") {
		p.desk.Next()
		left, err := p.calcExpr()
		if err != nil {
			return nil, err
		}
		if !p.desk.IsLabel(")") {
			return nil, fmt.Errorf("[ERROR](%d) Calc ')' not found", t.Line)
		}
		p.desk.Next() // ")"
		return left, nil
	}
	// simple value
	left, err := p.readValue1()
	if err != nil {
		return nil, err
	}
	return left, nil
}

func (p *Parser) readValue() (*node.Node, error) {
	if p.desk.IsLabel("=") {
		p.desk.Next()
	}
	n, err := p.calcExpr()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (p *Parser) readValue1() (*node.Node, error) {
	ct := p.desk.Peek()
	if p.desk.IsType(token.Number) {
		nn := node.NewNumber(ct.Label)
		p.desk.Next()
		return nn, nil
	} else if p.desk.IsType(token.Word) || p.desk.IsType(token.Macro) {
		wn := node.NewPushVariable(ct.Label)
		p.desk.Next()
		return wn, nil
	} else if p.desk.IsType(token.String) {
		sn := node.NewPushStr(ct.Label)
		p.desk.Next()
		return sn, nil
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
	case node.SetPitchBend:
		return node.NewSetPitchBend(no, opt), nil
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

func (p *Parser) readInt() (*node.Node, error) {
	errMsg := "[ERROR](%d) Invalid Int command : Int Name = Value"
	varName := p.desk.Next()
	if varName.Type != token.Word {
		return nil, fmt.Errorf(errMsg, varName.Line)
	}
	name := varName.Label
	var nodeInt *node.Node
	if p.desk.IsLabel("=") {
		p.desk.Next() // skip "="
		if !p.desk.IsType(token.Number) {
			return nil, fmt.Errorf(errMsg, varName.Line)
		}
		value, err := p.readValue()
		if err != nil {
			return nil, err
		}
		nodeInt = node.NewIntLet(varName.Label, value)
	} else {
		nodeInt = node.NewIntLet(varName.Label, node.NewNumber("0"))
	}
	p.Variable.SetIValue(name, 0) // temporary set variable
	return nodeInt, nil
}

func (p *Parser) readStr() (*node.Node, error) {
	errMsg := "[ERROR](%d) Invalid Str command : Str Name = {Value}"
	varName := p.desk.Next()
	if varName.Type != token.Word {
		return nil, fmt.Errorf(errMsg, varName.Line)
	}
	if !p.desk.IsLabel("=") {
		return nil, fmt.Errorf(errMsg, varName.Line)
	}
	var nodeStr *node.Node
	name := varName.Label
	if p.desk.IsLabel("=") {
		p.desk.Next() // skip "="
		nodeValue, err := p.readValue()
		if err != nil {
			return nil, fmt.Errorf(errMsg, varName.Line)
		}
		nodeStr = node.NewStrLet(name, nodeValue)
	} else {
		nodeStr = node.NewStrLet(name, node.NewPushStr(""))
	}
	p.Variable.SetSValue(name, "") // temporary set variable
	return nodeStr, nil
}

func (p *Parser) readMacro() (*node.Node, error) {
	errMsg := "[ERROR](%d) Invalid Macro command"
	// check name
	macroName := p.desk.Next()
	if macroName.Type != token.Macro {
		return nil, fmt.Errorf(errMsg, macroName.Line)
	}
	// call or define
	if p.desk.IsLabel("=") { // DEFINE MACRO
		p.desk.Next() // skip "="
		nodeValue, err := p.readValue()
		if err != nil {
			return nil, err
		}
		nodeStr := node.NewStrLet(macroName.Label, nodeValue)
		return nodeStr, nil
	} else {
		// Call macro
		callNode := node.NewStrEval(macroName.Label)
		return callNode, nil
	}
}

// Parse convert to AST
func Parse(tokens token.Tokens) (*node.Node, error) {
	return NewParser(tokens).Parse()
}
