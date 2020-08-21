package mmlparser

import (
	"github.com/kujirahand/sakuramml-go/song"
)

var moveNode *Node = nil

// Run : 実行
func Run(node *Node, so *song.Song) error {
	return runNode(node, so)
}

func runNode(node *Node, so *song.Song) error {
	err := node.Exec(node, so)
	if err != nil {
		return err
	}
	return nil
}

func runNodeList(node *Node, so *song.Song) error {
	so.Index = 0
	so.JumpTo = -1
	for so.Index < len(node.Children) {
		n := node.Children[so.Index]
		print("run:", n.ToString(0))
		err := n.Exec(n, so)
		if err != nil {
			return err
		}
		if so.JumpTo >= 0 {
			so.Index = so.JumpTo
			so.JumpTo = -1
			continue
		}
		so.Index++
	}
	return nil
}

func runNumber(node *Node, so *song.Song) error {
	num := node.Data.(ValueData).num
	so.PushValue(song.SNumber(num))
	return nil
}

func runNop(node *Node, song *song.Song) error {
	song.LastLineNo = node.Line.LineNo
	return nil
}

func runTone(node *Node, song *song.Song) error {
	data := node.Data.(ToneData)
	println("runTone=", data.Name, data.Flag)
	return nil
}

func runLoopBegin(node *Node, so *song.Song) error {
	cnt := 2
	if node.Children != nil || len(node.Children) > 0 {
		n := node.Children[0]
		n.Exec(n, so)
		cnt = so.PopIValue()
	}
	it := song.LoopItem{Count: cnt, Index: 0, Start: so.Index + 1, End: -1}
	so.PushLoop(&it)
	return nil
}

func runLoopEnd(node *Node, so *song.Song) error {
	it := so.PeekLoop()
	it.Index++
	if it.Count >= it.Index {
		so.PopLoop()
	} else {
		it.End = so.Index + 1
	}
	return nil
}

func runLoopBreak(node *Node, so *song.Song) error {
	it := so.PeekLoop()
	if (it.Count - 1) == it.Index {
		so.PopLoop()
		so.Index = it.End
	}
	return nil
}
