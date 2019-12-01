package song

import (
	"fmt"
	"sakuramml/track"
	"sakuramml/variable"
	"strconv"
)

// LoopItem struct
type LoopItem struct {
	Index     int
	Count     int
	BeginNode interface{} // Node
	EndNode   interface{} // Node
}

// EvalFunc type
type EvalStrFunc func (song *Song, src string) error

// Song is info of song, include tracks
type Song struct {
	Debug     bool
	Timebase  int
	Tempo     int
	TrackNo   int
	Stack     []interface{} // values stack
	Tracks    []*track.Track
	LoopStack []*LoopItem
	MoveNode  interface{} // Node
	Variable  *variable.Variable
	Eval      EvalStrFunc
}

// NewSong func
func NewSong() *Song {
	s := Song{}
	s.Debug = false
	s.Timebase = 96
	s.Tracks = []*track.Track{}
	// create default track
	for i := 0; i < 16; i++ {
		track := track.NewTrack(i, s.Timebase)
		s.Tracks = append(s.Tracks, track)
	}
	s.Tempo = 120 // default Tempo (but not write)
	s.TrackNo = 0
	s.Stack = make([]interface{}, 0, 256)
	s.MoveNode = nil
	s.Variable = variable.NewVariable()
	return &s
}

// PopStack func
func (song *Song) PopStack() interface{} {
	ilen := len(song.Stack)
	if ilen > 0 {
		iv := song.Stack[ilen-1]
		song.Stack = song.Stack[0 : ilen-1]
		return iv
	}
	return nil
}

// PushIValue func
func (song *Song) PushIValue(v int) {
	song.Stack = append(song.Stack, v)
}

// PopIValue func
func (song *Song) PopIValue() int {
	v := song.PopStack()
	switch v.(type) {
	case int:
		return v.(int)
	}
	return 0
}

// PushValue func
func (song *Song) PushValue(v *variable.Value) {
	if v.Type == variable.VTypeInt {
		song.PushIValue(v.IValue)
	} else if v.Type == variable.VTypeStr {
		song.PushSValue(v.SValue)
	} else {
		song.Stack = append(song.Stack, v)
	}
}

// PushSValue func
func (song *Song) PushSValue(v string) {
	song.Stack = append(song.Stack, v)
}

// PopSValue func
func (song *Song) PopSValue() string {
	v := song.PopStack()
	switch v.(type) {
	case string:
		return v.(string)
	}
	return ""
}

// PushLoop func
func (song *Song) PushLoop(item *LoopItem) {
	song.LoopStack = append(song.LoopStack, item)
}

// PopLoop func
func (song *Song) PopLoop() {
	ilen := len(song.LoopStack)
	if ilen > 0 {
		song.LoopStack = song.LoopStack[0 : ilen-1]
		return
	}
}

// PeekLoop func
func (song *Song) PeekLoop() *LoopItem {
	ilen := len(song.LoopStack)
	if ilen > 0 {
		return song.LoopStack[ilen-1]
	}
	return nil
}

// TimePtrToStr func
func (song *Song) TimePtrToStr(time int) string {
	l1 := song.Timebase * 4
	l4 := song.Timebase * 1
	meas := int(time / l1)
	beat := ((time - meas*l1) / l4)
	step := time % l4
	return fmt.Sprintf("%3d:%2d:%3d", meas, beat+1, step)
}

// StepToN lの逆引き
func (song *Song) StepToN(length int) string {
	tb := song.Timebase
	l1 := tb * 4
	if length > l1 {
		// 全音符より大きい場合
		nx := length / l1
		s := ""
		for i := 0; i < nx; i++ {
			s += "^"
		}
		res := "1" + s
		nm1 := length % l1
		if nm1 > 0 {
			res += "^" + song.StepToN(nm1)
		}
		return res
	}
	nd := l1 / length
	nm := l1 % length
	if nm == 0 {
		return strconv.Itoa(nd)
	}
	// 付点音符
	switch length {
	case int(float64(l1) / 2 * 1.5):
		return "2."
	case int(float64(l1) / 4 * 1.5):
		return "4."
	case int(float64(l1) / 8 * 1.5):
		return "8."
	case int(float64(l1) / 16 * 1.5):
		return "16."
	}
	return fmt.Sprintf("%%%d", length)
}

// CurTrack func
func (song *Song) CurTrack() *track.Track {
	for song.TrackNo >= len(song.Tracks) {
		tr := track.NewTrack(song.TrackNo%16, song.Timebase)
		song.Tracks = append(song.Tracks, tr)
	}
	return song.Tracks[song.TrackNo]
}

// ToString conver to string
func (song *Song) ToString() string {
	s := "Timebase=" + fmt.Sprintf("%d", song.Timebase) + "\n"
	for i := 0; i < len(song.Tracks); i++ {
		track := song.Tracks[i]
		if len(track.Events) == 0 {
			continue
		}
		s += fmt.Sprintf("+Track=%d\n", (i + 1))
		s += track.ToString()
	}
	return s
}
