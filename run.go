package sakuramml

import "fmt"

// Run : 実行
func SakuraRun(node *Node, so *Song) error {
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
		SakuraLog("run: " + n.ToString(0))
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
	numData := node.Data.(ValueData)
	so.PushSValue(numData.str)
	return nil
}

func runStr(node *Node, so *Song) error {
	numData := node.Data.(ValueData)
	so.PushSValue(numData.str)
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
	trk := song.CurTrack()
	toneData := node.Data.(ToneData)
	// rest ?
	noteNo := -1
	if toneData.Name == 'r' {
		noteNo = -1
	} else if toneData.Name == 'n' {
		toneData.NoteNo.Exec(toneData.NoteNo, song)
		noteNo = InRange(0, song.PopIValue(), 127)
	} else {
		// calc note
		o := trk.Octave
		if trk.OctaveOnce != 0 {
			o += trk.OctaveOnce
			trk.OctaveOnce = 0
		}
		noteNo = o*12 + noteToNo(toneData.Name)
		switch toneData.Flag {
		case "+", "#":
			noteNo += 1
		case "-":
			noteNo -= 1
		}
		noteNo = InRange(0, noteNo, 127)
	}
	// calc length
	length := trk.Length
	if toneData.Length != nil {
		toneData.Length.Exec(toneData.Length, song)
		length = song.PopStepValue()
	}
	gate := length
	if trk.QgateMode == "step" {
		gate = trk.Qgate
	} else {
		gate = int(float64(length) * float64(trk.Qgate) / 100.0)
	}
	// NoteOn
	if noteNo >= 0 {
		trk.AddNoteOn(trk.Time+trk.Timing, noteNo, trk.Velocity, gate)
	}
	trk.Time += length
	return nil
}

func runCommand(node *Node, song *Song) error {
	var v SValue = SStr("")
	trk := song.CurTrack()
	data := node.Data.(CommandData)
	// Get Command Value
	if data.Value != nil {
		err := data.Value.Exec(data.Value, song)
		if err == nil {
			v = song.PopSValue()
		}
	}

	switch string(data.Name) {
	case "v":
		trk.Velocity = v.ToInt()
	case "l":
		trk.Length = song.StrToStep(v.ToStr())
	case "q":
		trk.Qgate = v.ToInt()
	case "t":
		trk.Timing = v.ToInt()
	case "o":
		trk.Octave = InRange(0, v.ToInt(), 10)
	case "VOICE", "Voice", "@":
		// msb, lsb
		if data.Value2 != nil && data.Value3 != nil {
			data.Value2.Exec(data.Value2, song)
			msb := song.PopIValue()
			data.Value3.Exec(data.Value3, song)
			lsb := song.PopIValue()
			trk.AddCC(trk.Time-1, 0, msb)
			trk.AddCC(trk.Time-1, 0x20, lsb)
		}
		trk.AddProgramChange(trk.Time, v.ToInt())
	case "TR", "Track", "TRACK":
		song.TrackNo = v.ToInt()
	case "CH", "Channel", "CHANNEL":
		trk.Channel = InRange(1, v.ToInt(), 16) - 1
	case ">":
		trk.Octave = InRange(0, trk.Octave+1, 10)
	case "<":
		trk.Octave = InRange(0, trk.Octave-1, 10)
	case "`":
		trk.OctaveOnce += 1
	case "\"":
		trk.OctaveOnce -= 1
	case "y":
		if data.Value2 != nil {
			data.Value2.Exec(data.Value2, song)
			v2 := song.PopIValue()
			trk.AddCC(trk.Time, v.ToInt(), v2)
		}
	}
	// println("@@@runCommand=", string(data.Name), v.ToStr())
	return nil
}

func runLoopBegin(node *Node, so *Song) error {
	cnt := 2
	if node.Children != nil && len(node.Children) > 0 {
		tmp := so.Index
		runNodeList(node, so)
		s := so.PopSValue()
		cnt = s.ToInt()
		// SakuraLog(fmt.Sprintf("Loop=%s,%d", s.ToStr(), cnt))
		so.Index = tmp
	}
	it := LoopItem{Count: cnt, Index: 0, Start: so.Index + 1, End: -1}
	so.PushLoop(&it)
	return nil
}

func runLoopEnd(node *Node, so *Song) error {
	it := so.PeekLoop()
	it.End = so.Index + 1
	it.Index++
	if it.Count <= it.Index {
		so.PopLoop()
		so.JumpTo = it.End
	} else {
		so.JumpTo = it.Start
	}
	return nil
}

func runLoopBreak(node *Node, so *Song) error {
	it := so.PeekLoop()
	if (it.Count - 1) == it.Index {
		so.PopLoop()
		so.JumpTo = it.End
	}
	return nil
}

func runTime(node *Node, song *Song) error {
	// Calc Time (SakuraObj_time2step)
	// (ref) https://github.com/kujirahand/sakuramml-c/blob/68b62cbc101669211c511258ae1cf830616f238e/src/k_main.c#L473
	timeData := node.Data.(TimeData)
	timeData.v1.Exec(timeData.v1, song)
	timeData.v2.Exec(timeData.v2, song)
	timeData.v3.Exec(timeData.v3, song)
	tick := song.PopIValue()
	beat := song.PopIValue()
	mes := song.PopIValue()
	//
	base := song.Timebase * 4 / song.TimeSigDeno
	total := (mes-1)*(base*song.TimeSigFrac) + (beat-1)*base + tick
	song.CurTrack().Time = total
	return nil
}

func runTimeSig(node *Node, song *Song) error {
	timeData := node.Data.(TimeData)
	timeData.v1.Exec(timeData.v1, song)
	timeData.v2.Exec(timeData.v2, song)
	v2 := song.PopIValue()
	v1 := song.PopIValue()
	SakuraLog(fmt.Sprintf("TimeSig=%d,%d\n", v1, v2))
	song.TimeSigFrac = v1
	song.TimeSigDeno = v2
	return nil
}

func runLet(node *Node, song *Song) error {
	params := node.Data.(ParamsData)
	tag := params.tag
	name := params.name
	if tag == "INT" {
		params.v1.Exec(params.v1, song)
		iv := song.PopIValue()
		song.Variable.SetIValue(name, iv)
	} else if tag == "STR" {
		params.v1.Exec(params.v1, song)
		sv := song.PopSValue()
		song.Variable.SetSValue(name, sv.ToStr())
	}
	return nil
}

func runGetVar(node *Node, song *Song) error {
	params := node.Data.(ParamsData)
	name := params.name
	val := song.Variable.GetValue(name)
	// fmt.Printf("getvar %s = %s\n", name, val.ToString())
	song.PushSValue(val.ToString())
	return nil
}

func runPrint(node *Node, song *Song) error {
	params := node.Data.(ParamsData)
	params.v1.Exec(params.v1, song)
	sv := song.PopSValue()
	fmt.Println(sv)
	return nil
}
