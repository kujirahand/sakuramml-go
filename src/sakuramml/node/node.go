package node

import (
	"fmt"
	"sakuramml/song"
	"sakuramml/track"
	"sakuramml/utils"
	"strconv"
)

// NType type
type NType string

const (
	// Nop const
	Nop NType = "Nop"
	// Comment const
	Comment = "Comment"
	// NoteOn const
	NoteOn = "NoteOn"
	// SetTrack const
	SetTrack = "SetTrack"
	// SetTempo const
	SetTempo = "SetTempo"
	// SetOctave const
	SetOctave = "SetOctave"
	// SetQgate const
	SetQgate = "SetQgate"
	// SetVelocity const
	SetVelocity = "SetVelocity"
	// SetPC const
	SetPC = "@"
	// SetPitchBend const
	SetPitchBend = "SetPitchBend"
	// Number const
	Number = "Number"
	// Length const
	Length = "Length"
	// LengthDot const
	LengthDot = "Dot"
	// SetLength const
	SetLength = "SetLength"
	// GetTrackLength const
	GetTrackLength = "GetTrackLength"
	// CalcAdd const
	CalcAdd = "CalcAdd"
	// CalcMul const
	CalcMul = "CalcMul"
	// NLenToStep const
	NLenToStep = "NLenToStep"
	// LoopBegin const
	LoopBegin = "LoopBegin"
	// LoopEnd const
	LoopEnd = "LoopEnd"
	// LoopBreak const
	LoopBreak = "LoopBreak"
	// IntLet const
	IntLet = "IntLet"
	// StrLet const
	StrLet = "StrLet"
	// PushVariable const
	PushVariable = "PushVariable"
)

// Node struct
type Node struct {
	Type   NType
	Next   *Node
	Exec   func(n *Node, s *song.Song)
	IValue int
	SValue string
	NValue *Node
	ExData interface{}
}

func nodeToStringN(n *Node, level int) string {
	s := ""
	i := n
	for i != nil {
		// indent
		tab := ""
		for j := 0; j < level; j++ {
			tab += "|  "
		}
		///fmt.Printf(tab+"%d %v\n", level, *i)
		params := ""
		switch i.Type {
		case NoteOn:
			params = i.SValue
		case Number:
			params = fmt.Sprintf("%d", i.IValue)
		case PushVariable:
			params = fmt.Sprintf("%s", i.SValue)
		}
		s += tab + string(i.Type) + " " + params + "\n"
		if i.NValue != nil {
			s += nodeToStringN(i.NValue, level+1)
		}
		i = i.Next
	}
	return s
}

// ToStringAll func
func (n *Node) ToStringAll() string {
	return nodeToStringN(n, 0)
}

// ToStringAllName func
func (n *Node) ToStringAllName(delimiter string) string {
	s := ""
	nTop := n
	nCur := nTop
	for nCur != nil {
		s += string(nCur.Type)
		if nCur.Next != nil {
			s += delimiter
		}
		nCur = nCur.Next
	}
	return s
}

// ExDataNode struct
type ExDataNode struct {
	Value *Node
}

// ExDataNoteOn struct
type ExDataNoteOn struct {
	NoteShift   int
	Length      *Node
	Qgate       *Node
	QgateOpt    rune
	QgateMode   string
	Velocity    *Node
	VelocityOpt rune
	NoteNo      *Node
	NoteNoOpt   rune
}

// NewNode func
func NewNode(nodeType NType) *Node {
	n := Node{Type: nodeType, Exec: execNone}
	n.Next = nil
	n.NValue = nil
	return &n
}

func execNone(n *Node, s *song.Song) {
	err := fmt.Errorf("Exec failed, not implemented : %v", *n)
	panic(err)
}

// NewNop func
func NewNop() *Node {
	t := NewNode(Nop)
	t.Exec = execNop
	return t
}
func execNop(n *Node, s *song.Song) {
	// nop
}

// NewComment func
func NewComment(text string) *Node {
	t := NewNode(Comment)
	t.Exec = execComment
	t.SValue = text
	return t
}
func execComment(n *Node, s *song.Song) {
	tr := s.CurTrack()
	tb := []byte(n.SValue)
	if len(tb) > 255 {
		tb = tb[0:255]
	}
	tr.AddMeta(tr.Time, 0x01, tb)
}

// NewNoteOn func (NoteOn and Rest)
func NewNoteOn(note string, ex *ExDataNoteOn) *Node {
	// detect note no
	notemap := map[string]int{
		"c": 0, "d": 2, "e": 4, "f": 5, "g": 7, "a": 9, "b": 11, "r": -1, "n": -1,
	}
	// new
	n := NewNode(NoteOn)
	n.Exec = execNoteOn
	n.SValue = note
	n.IValue = notemap[note]
	n.ExData = ex
	return n
}

func execNoteOn(n *Node, s *song.Song) {
	tr := s.CurTrack()
	noteno := 0
	length := tr.Length
	qgate := tr.Qgate
	qgatemode := tr.QgateMode
	qgateAdd := 0
	velocity := tr.Velocity
	// Temporary change?
	ex := n.ExData.(*ExDataNoteOn)
	if ex.Length != nil {
		ex.Length.Exec(ex.Length, s)
		length = s.PopIValue()
	}
	if ex.Qgate != nil {
		ex.Qgate.Exec(ex.Qgate, s)
		qv := s.PopIValue()
		if ex.QgateOpt == rune('+') || ex.QgateOpt == rune('-') {
			if ex.QgateMode == track.QgateModeStep {
				qgateAdd = qv
			} else {
				qgateAdd = int(float64(length) * float64(qv) / 100)
			}
			if ex.QgateOpt == rune('-') {
				qgateAdd *= -1
			}
		} else {
			qgate = calcFlagValue(qgate, s.PopIValue(), string(ex.QgateOpt))
			qgatemode = ex.QgateMode
		}
	}
	if ex.Velocity != nil {
		ex.Velocity.Exec(ex.Velocity, s)
		velocity = calcFlagValue(velocity, s.PopIValue(), string(ex.VelocityOpt))
	}
	// calc
	qlen := qgate
	if qgatemode == track.QgateModeRate {
		qlen = int(float64(length) * float64(qgate) / 100)
	}
	qlen += qgateAdd
	// rest or note
	if n.SValue == "r" {
		if s.Debug {
			nls := s.StepToN(length)
			fmt.Printf("- Time(%s) l%-2s r \n", s.TimePtrToStr(tr.Time), nls)
		}
	} else {
		if n.SValue == "n" {
			ex.NoteNo.Exec(ex.NoteNo, s)
			noteno = s.PopIValue()
		} else {
			// calc note shift(# or flat)
			noteno = tr.Octave*12 + n.IValue + ex.NoteShift
		}
		if s.Debug {
			notemap2 := []string{"c", "c#", "d", "d#", "e", "f", "f#", "g", "g#", "a", "a#", "b"}
			nls := s.StepToN(length)
			fmt.Printf(
				"- Time(%s) TR=%-2d l%-2s o%d v%-3d q%%%-3d %-3s \n",
				s.TimePtrToStr(tr.Time), s.TrackNo, nls, int(noteno/12), velocity, qlen, notemap2[noteno%12])
		}
		tr.AddNoteOn(tr.Time, noteno, velocity, qlen)
	}
	tr.Time += length
}

// NewNumber func
func NewNumber(s string) *Node {
	if s == "" {
		return nil
	}
	base := 10
	if len(s) > 2 && s[0:2] == "0x" {
		base = 16
	}
	iv, _ := strconv.ParseInt(s, base, 0)
	n := NewNode(Number)
	n.Exec = execPushIValue
	n.IValue = int(iv)
	return n
}

// NewNumberInt func
func NewNumberInt(no int) *Node {
	n := NewNode(Number)
	n.Exec = execPushIValue
	n.IValue = no
	return n
}

func execPushIValue(n *Node, s *song.Song) {
	s.PushIValue(n.IValue)
}

// NewPushVariable func
func NewPushVariable(word string) *Node {
	if word == "" {
		return nil
	}
	n := NewNode(PushVariable)
	n.SValue = word
	n.Exec = execPushVariable
	return n
}

func execPushVariable(n *Node, s *song.Song) {
	word := n.SValue
	value := s.Variable.GetValue(word)
	s.PushValue(value)
}


// NewSetTrack func
func NewSetTrack(v *Node, opt string) *Node {
	n := NewNode(SetTrack)
	n.Exec = execSetTrack
	n.SValue = opt
	n.NValue = v
	return n
}
func execSetTrack(n *Node, s *song.Song) {
	// get track no
	n.NValue.Exec(n.NValue, s)
	// set new value
	s.TrackNo = calcFlagValue(s.TrackNo, s.PopIValue(), n.SValue)
}

func calcFlagValue(cur, no int, opt string) int {
	res := cur
	switch opt {
	case "+":
		res += no
	case "-":
		res -= no
	case "++":
		res++
	case "--":
		res--
	default:
		res = no
	}
	if res < 0 {
		res = 0
	}
	return res
}

// NewSetOctave func
func NewSetOctave(v *Node, opt string) *Node {
	n := NewNode(SetOctave)
	n.Exec = execSetOctave
	n.SValue = opt
	n.NValue = v
	return n
}

func execSetOctave(n *Node, s *song.Song) {
	tr := s.CurTrack()
	no := 0
	if n.NValue != nil {
		n.NValue.Exec(n.NValue, s)
		no = s.PopIValue()
	}
	tr.Octave = calcFlagValue(tr.Octave, no, n.SValue)
	if tr.Octave > 10 {
		tr.Octave = 10
	}
}

// NewSetVelocity func
func NewSetVelocity(v *Node, opt string) *Node {
	n := NewNode(SetVelocity)
	n.Exec = execSetVelocity
	n.SValue = opt
	n.NValue = v
	return n
}

func execSetVelocity(n *Node, s *song.Song) {
	n.NValue.Exec(n.NValue, s)
	tr := s.CurTrack()
	tr.Velocity = calcFlagValue(tr.Velocity, s.PopIValue(), n.SValue)
	if tr.Velocity > 127 {
		tr.Velocity = 127
	}
}

// NewSetQgate func
func NewSetQgate(v *Node, opt string) *Node {
	n := NewNode(SetQgate)
	n.Exec = execSetQgate
	n.SValue = opt
	n.NValue = v
	return n
}

func execSetQgate(n *Node, s *song.Song) {
	n.NValue.Exec(n.NValue, s)
	tr := s.CurTrack()
	opt := n.SValue
	// set Qgate
	if len(opt) > 0 && opt[0] == '%' {
		// Direct Value
		opt = opt[1:]
		tr.Qgate = calcFlagValue(tr.Qgate, s.PopIValue(), opt)
		tr.QgateMode = track.QgateModeStep
	} else {
		// Percent Value
		tr.Qgate = calcFlagValue(tr.Qgate, s.PopIValue(), opt)
		tr.QgateMode = track.QgateModeRate
	}
	if tr.Qgate < 1 {
		tr.Qgate = 1
	}
}

// NewSetTempo func
func NewSetTempo(v *Node, opt string) *Node {
	n := NewNode(SetTempo)
	n.Exec = execSetTempo
	n.SValue = opt
	n.NValue = v
	return n
}

func execSetTempo(n *Node, s *song.Song) {
	n.NValue.Exec(n.NValue, s)
	s.Tempo = calcFlagValue(s.Tempo, s.PopIValue(), n.SValue)
	s.Tempo = utils.InRange(10, s.Tempo, 1500)
	trk := s.CurTrack()
	trk.AddTempo(trk.Time, s.Tempo)
}

// NewSetPitchBend func
func NewSetPitchBend(v *Node, opt string) *Node {
	n := NewNode(SetPitchBend)
	n.Exec = execSetPitchBend
	n.SValue = opt
	n.NValue = v
	return n
}

func execSetPitchBend(n *Node, s *song.Song) {
	tr := s.CurTrack()
	n.NValue.Exec(n.NValue, s)
	opt := n.SValue
	// PitchBend Mode
	if len(opt) > 0 && opt[0] == '%' { // normal mode
		opt = opt[1:]
		tr.PitchBend = calcFlagValue(tr.PitchBend, s.PopIValue(), opt)
		tr.PitchBend = utils.InRange(-8192, tr.PitchBend, 8191)
		tr.AddPitchBend(tr.Time, tr.PitchBend)
	} else { // easy mode
		pb := calcFlagValue(tr.PitchBend, s.PopIValue(), opt)
		pb = utils.InRange(0, pb, 127)
		tr.AddPitchBendEasy(tr.Time, tr.PitchBend)
	}
}

// ExDataPC for SetPC
type ExDataPC struct {
	MSB int
	LSB int
}

// NewSetPC func
func NewSetPC(v *Node, msb, lsb int) *Node {
	n := NewNode(SetPC)
	n.Exec = execSetPC
	n.NValue = v
	n.ExData = ExDataPC{MSB: msb, LSB: lsb}
	return n
}

func execSetPC(n *Node, s *song.Song) {
	// track
	tr := s.CurTrack()
	// value
	n.NValue.Exec(n.NValue, s)
	no := utils.MidiRange(s.PopIValue() - 1)
	// write msb lsb
	ex := n.ExData.(ExDataPC)
	if ex.MSB >= 0 {
		tr.AddCC(tr.Time, 0, ex.MSB)
	}
	if ex.LSB >= 0 {
		tr.AddCC(tr.Time, 32, ex.LSB)
	}
	// write pc
	tr.AddProgramChange(tr.Time, no)
	//
	if s.Debug {
		fmt.Printf(
			"- Time(%s) TR=%-2d @%d,%d, %d\n",
			s.TimePtrToStr(tr.Time), s.TrackNo, no+1, ex.MSB, ex.LSB)
	}
}

// NewLength func
func NewLength() *Node {
	n := NewNode(Length)
	n.Exec = execLength
	return n
}

func execLength(n *Node, s *song.Song) {
	// calc length
	length := 0
	nvalue := n.NValue
	i := 0
	for nvalue != nil {
		if nvalue.Type == Nop {
			nvalue = nvalue.Next
			continue
		}
		// if s.Debug { fmt.Printf("%d, %s\n", i, nvalue.Type) }
		nvalue.Exec(nvalue, s)
		iv := s.PopIValue()
		length += iv
		nvalue = nvalue.Next
		i++
		if i > 10 {
			break
		}
	}
	s.PushIValue(length)
}

// NewSetLength func
func NewSetLength(lenNode *Node) *Node {
	n := NewNode(SetLength)
	n.NValue = lenNode
	n.Exec = execSetLength
	return n
}

func execSetLength(n *Node, s *song.Song) {
	n.NValue.Exec(n, s)
	ilen := s.PopIValue()
	// println("execSetLength=", ilen)
	s.CurTrack().Length = ilen
}

// NewGetTrackLength func
func NewGetTrackLength() *Node {
	n := NewNode(GetTrackLength)
	n.Exec = execGetTrackLength
	return n
}

func execGetTrackLength(n *Node, s *song.Song) {
	s.PushIValue(s.CurTrack().Length)
}

// NewLengthDot func
func NewLengthDot(nLen *Node) *Node {
	n := NewNode(LengthDot)
	n.Exec = execLenDot
	n.NValue = nLen
	n.ExData = float64(1.5)
	return n
}

func execLenDot(n *Node, s *song.Song) {
	rate := n.ExData.(float64)
	// get number
	n.NValue.Exec(n.NValue, s)
	iv := s.PopIValue()
	// calc len
	vv := int(float64(iv) * rate)
	s.PushIValue(vv)
	// println("dot=", iv, rate, vv)
}

// NewCalcAdd func
func NewCalcAdd(lnode, rnode *Node) *Node {
	n := NewNode(CalcAdd)
	n.Exec = execCalcAdd
	n.ExData = []*Node{lnode, rnode}
	return n
}

func execCalcAdd(n *Node, s *song.Song) {
	ex := n.ExData.([]*Node)
	lnode, rnode := ex[0], ex[1]
	rnode.Exec(n, s)
	rvalue := s.PopIValue()
	lnode.Exec(n, s)
	lvalue := s.PopIValue()
	vv := rvalue + lvalue
	s.PushIValue(vv)
}

// NewCalcMul func
func NewCalcMul(lnode, rnode *Node) *Node {
	n := NewNode(CalcMul)
	n.Exec = execCalcMul
	n.ExData = []*Node{lnode, rnode}
	return n
}

func execCalcMul(n *Node, s *song.Song) {
	ex := n.ExData.([]*Node)
	lnode, rnode := ex[0], ex[1]
	rnode.Exec(n, s)
	rvalue := s.PopIValue()
	lnode.Exec(n, s)
	lvalue := s.PopIValue()
	vv := rvalue * lvalue
	s.PushIValue(vv)
}

// NewNLenToStep func
func NewNLenToStep(valueNode *Node) *Node {
	n := NewNode(NLenToStep)
	n.Exec = execNLenToStep
	n.NValue = valueNode
	return n
}

func execNLenToStep(n *Node, s *song.Song) {
	// get n value
	nValue := n.NValue
	nValue.Exec(nValue, s)
	v := s.PopIValue()
	// convert to step
	vStep := int((4.0 / float64(v)) * float64(s.Timebase))
	s.PushIValue(vStep)
}

// NewLoopBegin func
func NewLoopBegin(loopValue *Node) *Node {
	n := NewNode(LoopBegin)
	n.Exec = execLoopBegin
	n.NValue = loopValue
	return n
}

func execLoopBegin(n *Node, s *song.Song) {
	loopValue := n.NValue
	loopValue.Exec(loopValue, s)
	loopCount := s.PopIValue()
	// Search LoopEndPoint
	var endPoint *Node = nil
	cur := n.Next
	for cur != nil {
		if cur.Type == LoopEnd {
			endPoint = cur
			break
		}
		cur = cur.Next
	}
	// loop item
	loop := song.LoopItem{
		Count:     loopCount,
		Index:     0,
		BeginNode: n.Next,
		EndNode:   endPoint.Next,
	}
	s.PushLoop(&loop)
}

// NewLoopEnd func
func NewLoopEnd() *Node {
	n := NewNode(LoopEnd)
	n.Exec = execLoopEnd
	return n
}

func execLoopEnd(n *Node, s *song.Song) {
	loop := s.PeekLoop()
	loop.Index++
	if loop.Index == loop.Count {
		s.PopLoop()
		return
	}
	// back to begin node
	s.MoveNode = loop.BeginNode
}

// NewLoopBreak func
func NewLoopBreak() *Node {
	n := NewNode(LoopBreak)
	n.Exec = execLoopBreak
	return n
}

func execLoopBreak(n *Node, s *song.Song) {
	loop := s.PeekLoop()
	// last one time?
	if loop.Index == loop.Count-1 {
		// go to last
		s.MoveNode = loop.EndNode
	}
}

func NewIntLet(name string, value int) *Node {
	n := NewNode(IntLet)
	n.SValue = name
	n.IValue = value
	n.Exec = execIntLet
	return n
}

func execIntLet(n *Node, s *song.Song) {
	varName := n.SValue
	s.Variable.SetIValue(varName, n.IValue)
}

