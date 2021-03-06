package mmlparser

import (
	"fmt"

	"github.com/kujirahand/sakuramml-go/song"
)

const (
	// Nop : Nop
	Nop = iota
	// NodeList : NodeList
	NodeList
	// NodeTone : cdefgab
	NodeTone
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
)

var nodeTypeMap map[int]string = map[int]string{
	Nop:           "Nop",
	NodeList:      "NodeList",
	NodeTone:      "NodeTone",
	NodeEOL:       "NodeEOL",
	NodeLoopBegin: "NodeLoopBegin",
	NodeLoopEnd:   "NodeLoopEnd",
	NodeLoopBreak: "NodeLoopBreak",
	NodeNumber:    "NodeNumber",
}

// ExecFunc func
type ExecFunc func(n *Node, s *song.Song) error

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
	Name rune
	Flag string
}

func (p ToneData) toString() string {
	return string(p.Name) + p.Flag
}

// NewToneNode : tone node
func NewToneNode(tok Token, flag string) *Node {
	node := NewNode(NodeTone)
	node.Data = ToneData{
		Name: rune(tok.label[0]),
		Flag: flag,
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
	if expr == nil {
		expr = NewNumberNode(Token{label: "2", value: 2})
	}
	n.Children = []*Node{expr}
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
