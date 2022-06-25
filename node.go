package sakuramml

import (
	"fmt"
)

const (
	// Nop : Nop
	Nop = iota
	// NodeList : NodeList
	NodeList
	// NodeTone : cdefgab
	NodeTone
	// NodeCommand : vloq
	NodeCommand
	// NodeEOL : End of Line
	NodeEOL
	// NodeLoopBegin : NodeLoopBegin
	NodeLoopBegin
	// NodeLoopEnd : NodeLoopEnd
	NodeLoopEnd
	// NodeLoopBreak : NodeLoopBreak
	NodeLoopBreak
	// NodeNumber : number
	NodeNumber
	// NodeTime : Time
	NodeTime
	// NodeTimeSig : TimeSignature
	NodeTimeSig
)

var nodeTypeMap map[int]string = map[int]string{
	Nop:           "Nop",
	NodeList:      "NodeList",
	NodeTone:      "NodeTone",
	NodeCommand:   "NodeCommand",
	NodeEOL:       "NodeEOL",
	NodeLoopBegin: "NodeLoopBegin",
	NodeLoopEnd:   "NodeLoopEnd",
	NodeLoopBreak: "NodeLoopBreak",
	NodeNumber:    "NodeNumber",
	NodeTime:      "NodeTime",
	NodeTimeSig:   "NodeTimeSig",
}

// ExecFunc func
type ExecFunc func(n *Node, s *Song) error

// LineInfo : ソースコードの場所
type LineInfo struct {
	LineNo int
	FileNo int
}

// NodeData : Node Data
type NodeData interface {
	toString() string
}

// Node : Node
type Node struct {
	Type     int
	Children []*Node
	Line     LineInfo
	Data     NodeData
	Exec     ExecFunc
}

// EmptyData ... empty
type EmptyData struct {
}

func (p EmptyData) toString() string {
	return ""
}

// ValueData : value data
type ValueData struct {
	num float64
	str string
}

func (p ValueData) toString() string {
	return p.str
}

// ToneData : Data
type ToneData struct {
	Name   rune
	Flag   string
	Length *Node
}

func (p ToneData) toString() string {
	return string(p.Name) + p.Flag
}

// NewToneNode : tone node
func NewToneNode(tok Token, flag string, len *Node) *Node {
	node := NewNode(NodeTone)
	node.Data = ToneData{
		Name:   rune(tok.label[0]),
		Flag:   flag,
		Length: len,
	}
	node.Exec = runTone
	return node
}

// ToString : debug
func (p *Node) ToString(level int) string {
	res := ""
	for i := 0; i < level; i++ {
		res += "|--"
	}
	name, ok := nodeTypeMap[p.Type]
	if !ok {
		name = "???(node.goで定義が必要)???"
	}
	res += name + ":" + p.Data.toString() + ":" + fmt.Sprintf("(%d)", p.Line.LineNo) + "\n"
	// children
	for _, n := range p.Children {
		res += n.ToString(level + 1)
	}
	return res
}

// CommandData: Data
type CommandData struct {
	Name  string
	Value *Node
}

func (p CommandData) toString() string {
	return string(p.Name) + p.Value.ToString(0)
}

// NewToneNode : tone node
func NewCommandNode(tok Token, name string, val *Node) *Node {
	if name == "WORD" {
		name = tok.label
	}
	data := CommandData{Name: name, Value: val}
	node := NewNode(NodeCommand)
	node.Data = data
	node.Exec = runCommand
	return node
}

// CurLine : 現在パース中のライン情報
var CurLine LineInfo = LineInfo{}

// NewNode : new node
func NewNode(nodeType int) *Node {

	n := Node{
		Type: nodeType,
		Line: CurLine,
		Data: EmptyData{},
		Exec: runNop,
	}

	if nodeType == NodeList {
		n.Exec = runNodeList
	}

	return &n
}

// NewNumberNode : NewNumberNode
func NewNumberNode(t Token) *Node {
	n := NewNode(NodeNumber)
	n.Data = ValueData{
		num: t.value,
		str: t.label,
	}
	n.Exec = runNumber
	return n
}

// NewLoopNodeBegin : loop begin
func NewLoopNodeBegin(t Token, expr *Node) *Node {
	n := NewNode(NodeLoopBegin)
	if expr != nil { // it has loop counter
		n.Children = []*Node{expr}
	}
	n.Exec = runLoopBegin
	return n
}

// NewLoopNodeEnd : end
func NewLoopNodeEnd(t Token) *Node {
	n := NewNode(NodeLoopEnd)
	n.Line = t.lineinfo
	n.Exec = runLoopEnd
	return n
}

// NewLoopNodeBreak : break
func NewLoopNodeBreak(t Token) *Node {
	n := NewNode(NodeLoopBreak)
	n.Line = t.lineinfo
	n.Exec = runLoopBreak
	return n
}

// AppendChildNode : append child node
func AppendChildNode(parent *Node, child *Node) *Node {
	parent.Children = append(parent.Children, child)
	return parent
}

// TimeData : TimeData
type TimeData struct {
	mode string
	v1   *Node
	v2   *Node
	v3   *Node
}

func (p TimeData) toString() string {
	return p.mode + "(?:?:?)"
}

// NewToneNode : tone node
func NewTimeNode(tok Token, v1 *Node, v2 *Node, v3 *Node) *Node {
	node := NewNode(NodeTime)
	node.Data = TimeData{
		mode: "Time",
		v1:   v1,
		v2:   v2,
		v3:   v3,
	}
	node.Exec = runTime
	return node
}

// NewToneNode : tone node
func NewTimeSigNode(tok Token, v1 *Node, v2 *Node) *Node {
	node := NewNode(NodeTimeSig)
	node.Data = TimeData{
		mode: "TimeSig",
		v1:   v1,
		v2:   v2,
		v3:   nil,
	}
	node.Exec = runTimeSig
	return node
}
