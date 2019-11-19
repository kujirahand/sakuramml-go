package node

import (
	"fmt"
	"sakuramml/song"
	"strconv"
)

const (
	// Nop const
	Nop = "Nop"
	// NoteOn const
	NoteOn = "NoteOn"
	// Track const
	Track = "Track"
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
)

// NType type
type NType string

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
		s += tab + string(i.Type) + "\n"
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

// ExDataNode strcut
type ExDataNode struct {
	Value *Node
}

// ExDataNoteOn struct
type ExDataNoteOn struct {
	Length *Node
}

// NewNode func
func NewNode(nodeType NType) *Node {
	n := Node{Type: nodeType, Exec: execNop}
	n.Next = nil
	n.NValue = nil
	return &n
}

// NewNop func
func NewNop() *Node {
	return NewNode(Nop)
}
func execNop(n *Node, s *song.Song) {
	// nop
}

// NewNoteOn func
func NewNoteOn(note string, ex *ExDataNoteOn) *Node {
	n := NewNode(NoteOn)
	n.Exec = execNoteOn
	n.SValue = note
	n.ExData = ex
	return n
}

func execNoteOn(n *Node, s *song.Song) {
	track := s.CurTrack()
	noteno := 0
	velocity := track.Velocity
	qgate := track.Qgate
	length := track.Length
	// Temporary change?
	ex := n.ExData.(*ExDataNoteOn)
	if ex.Length != nil {
		ex.Length.Exec(ex.Length, s)
		length = s.PopIValue()
	}
	if n.SValue == "n" {
		// todo "n"
	} else {
		notemap := map[string]int{"c": 0, "d": 2, "e": 4, "f": 5, "g": 7, "a": 9, "b": 11}
		noteno = track.Octave*12 + notemap[n.SValue]
		if s.Debug {
			print("- note:", noteno, ",%", length, "\n")
		}
	}
	track.AddNoteOn(track.Time, noteno, velocity, qgate)
	track.Time += length
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

func execPushIValue(n *Node, s *song.Song) {
	s.PushIValue(n.IValue)
}

// NewTrack func
func NewTrack(v *Node) *Node {
	n := NewNode(Track)
	n.Exec = execTrack
	n.ExData = ExDataNode{Value: v}
	return n
}

func execTrack(n *Node, s *song.Song) {
	// get track no
	ex := n.ExData.(ExDataNode)
	ex.Value.Exec(n, s)
	// Change Current TrackNo
	s.TrackNo = s.PopIValue()
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
		if s.Debug {
			fmt.Printf("%d, %s\n", i, nvalue.Type)
		}
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
