package node

import (
	"sakuramml/song"
	"strconv"
)

const (
	// TypeNop const
	TypeNop = "Nop"
	// TypeNoteOn const
	TypeNoteOn = "NoteOn"
	// TypeTrack const
	TypeTrack = "Track"
	// TypeNumber const
	TypeNumber = "Number"
)

// Node struct
type Node struct {
	Type   string
	Next   *Node
	Exec   func(n *Node, s *song.Song)
	IValue int
	SValue string
	ExData interface{}
}

// ExDataNode strcut
type ExDataNode struct {
	Value *Node
}

// NewNop func
func NewNop() *Node {
	n := Node{Type: TypeNop, Next: nil, Exec: execNop}
	return &n
}
func execNop(n *Node, s *song.Song) {
	// nop
}

// NewNoteOn func
func NewNoteOn(note string) *Node {
	n := Node{Type: TypeNoteOn, Next: nil, Exec: execNoteOn}
	n.SValue = note
	return &n
}

func execNoteOn(n *Node, s *song.Song) {
	track := s.CurTrack()
	noteno := 0
	velocity := track.Velocity
	qgate := track.Qgate
	length := track.Length
	if n.SValue == "n" {
		// todo "n"
	} else {
		notemap := map[string]int{"c": 0, "d": 2, "e": 4, "f": 5, "g": 7, "a": 9, "b": 11}
		noteno = track.Octave*12 + notemap[n.SValue]
		print("octave:", track.Octave, "\n")
		print("note:", noteno, "\n")
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
	n := Node{Type: TypeNumber, Next: nil, Exec: execPushIValue}
	n.IValue = int(iv)
	return &n
}

func execPushIValue(n *Node, s *song.Song) {
	s.PushIValue(n.IValue)
}

// NewTrack func
func NewTrack(v *Node) *Node {
	n := Node{Type: TypeTrack, Exec: execTrack}
	n.ExData = ExDataNode{Value: v}
	return &n
}

func execTrack(n *Node, s *song.Song) {
	// get track no
	ex := n.ExData.(ExDataNode)
	ex.Value.Exec(n, s)
	// Change Current TrackNo
	s.TrackNo = s.PopIValue()
}
