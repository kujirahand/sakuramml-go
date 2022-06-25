package sakuramml

import "log"

// Run : 実行
func Run(node *Node, so *Song) error {
	return runNode(node, so)
}

func runNode(node *Node, so *Song) error {
	err := node.Exec(node, so)
	if err != nil {
		return err
	}
	return nil
}

func runNodeList(node *Node, so *Song) error {
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

func runNumber(node *Node, so *Song) error {
	num := node.Data.(ValueData).num
	so.PushValue(SNumber(num))
	return nil
}

func runNop(node *Node, song *Song) error {
	song.LastLineNo = node.Line.LineNo
	return nil
}

func noteToNo(n rune) int {
	switch n {
	case 'c':
		return 0
	case 'd':
		return 2
	case 'e':
		return 4
	case 'f':
		return 5
	case 'g':
		return 7
	case 'a':
		return 9
	case 'b':
		return 11
	}
	return -1
}

func runTone(node *Node, song *Song) error {
	toneData := node.Data.(ToneData)
	println("runTone=", toneData.Name, toneData.Flag)
	trk := song.CurTrack()
	// calc note
	nn := trk.Octave*12 + noteToNo(toneData.Name)
	switch toneData.Flag {
	case "+":
		nn += 1
	case "-":
		nn -= 1
	}
	// calc length
	length := trk.Length
	gate := length
	if trk.QgateMode == "step" {
		gate = trk.Qgate
	} else {
		gate = int(float64(length) * float64(trk.Qgate) / 100.0)
	}
	trk.AddNoteOn(trk.Time, nn, trk.Velocity, gate)
	trk.Time += length
	return nil
}

func runCommand(node *Node, song *Song) error {
	v := 0
	trk := song.CurTrack()
	data := node.Data.(CommandData)
	err := data.Value.Exec(data.Value, song)
	if err != nil {
		switch data.Name {
		case 'v':
			v = 100
		case 'l':
			v = trk.Length
		case 'o':
			v = 5
		case 'q':
			v = 90
		}
	} else {
		log.Printf("cmd=%d", song.Stack[len(song.Stack)-1])
		v = song.PopIValue()
	}
	log.Printf("cmd=%d", v)
	switch data.Name {
	case 'v':
		trk.Velocity = v
	case 'l':
		trk.Length = song.NToStep(v)
	case 'q':
		trk.Qgate = v
	case 'o':
		trk.Octave = v
	}
	println("runCommand=", string(data.Name), v)
	return nil
}

func runLoopBegin(node *Node, so *Song) error {
	cnt := 2
	if node.Children != nil || len(node.Children) > 0 {
		n := node.Children[0]
		n.Exec(n, so)
		cnt = so.PopIValue()
	}
	it := LoopItem{Count: cnt, Index: 0, Start: so.Index + 1, End: -1}
	so.PushLoop(&it)
	return nil
}

func runLoopEnd(node *Node, so *Song) error {
	it := so.PeekLoop()
	it.Index++
	if it.Count >= it.Index {
		so.PopLoop()
	} else {
		it.End = so.Index + 1
	}
	return nil
}

func runLoopBreak(node *Node, so *Song) error {
	it := so.PeekLoop()
	if (it.Count - 1) == it.Index {
		so.PopLoop()
		so.Index = it.End
	}
	return nil
}
