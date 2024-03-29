package sakuramml

import (
	"fmt"
	"strconv"
)

// LoopItem struct
type LoopItem struct {
	Index int
	Count int
	Start int
	End   int
}

// EvalStrFunc type
type EvalStrFunc func(song *Song, src string) error

// Song is info of song, include tracks
type Song struct {
	Debug       bool
	Timebase    int
	TimeSigFrac int // 分子
	TimeSigDeno int // 分母
	Tempo       int
	TrackNo     int
	Stack       []SValue
	Tracks      []*Track
	LoopStack   []*LoopItem
	Variable    *Variable
	Eval        EvalStrFunc
	LastLineNo  int
	Index       int
	JumpTo      int
}

// NewSong func
func NewSong() *Song {
	s := Song{}
	s.Debug = false
	s.Timebase = 96
	s.TimeSigDeno = 4 //
	s.TimeSigFrac = 4 //
	s.Tracks = []*Track{}
	// create default track
	for i := 0; i < 16; i++ {
		track := NewTrack(i, s.Timebase)
		s.Tracks = append(s.Tracks, track)
	}
	s.Tempo = 120 // default Tempo (but not write)
	s.TrackNo = 0
	s.Stack = []SValue{}
	s.Variable = NewVariable()
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

// PopIValue func
func (song *Song) PopIValue() int {
	iv := song.PopStack()
	sv := iv.(SValue)
	return sv.ToInt()
}

// PopSValue func
func (song *Song) PopSValue() SValue {
	v := song.PopStack()
	return v.(SValue)
}

// PopStepValue func
func (song *Song) PopStepValue() int {
	iv := song.PopStack()
	sv := iv.(SValue).ToStr()
	return song.StrToStep(sv)
}

// PushIValue func
func (song *Song) PushIValue(v int) {
	song.Stack = append(song.Stack, SNumber(v))
}

// PushValue func
func (song *Song) PushValue(v SValue) {
	song.Stack = append(song.Stack, v)
}

// PushSValue func
func (song *Song) PushSValue(v string) {
	song.Stack = append(song.Stack, SStr(v))
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

func (song *Song) StrToStep(s string) int {
	trk := song.CurTrack()
	defLen := trk.Length
	if s == "" {
		return defLen
	}
	total := 0
	sl := newSLexer(s, 0)
	for !sl.isEOF() {
		num := defLen
		c := sl.peek()
		if c == '%' {
			sl.next() // skip %
			num = sl.readInt(defLen)
		} else if isDigit(c) {
			n := sl.readInt(0)
			num = song.NToStep(n)
			countDot := 0
			for sl.peek() == '.' {
				countDot += 1
				sl.next()
			}
			if countDot > 0 {
				switch countDot {
				case 1:
					num = int(float64(num) * 1.5)
				case 2:
					num = int(float64(num) * (1.0 + 0.5 + 0.25))
				case 3:
					num = int(float64(num) * (1.0 + 0.5 + 0.25 + 0.12))
				case 4:
					num = int(float64(num) * (1.0 + 0.5 + 0.25 + 0.12 + 0.6))
				default:
					num = int(float64(num) * (1.0 + 0.5 + 0.25 + 0.12 + 0.6 + 0.3))
				}
			}

		}
		total += num
		if sl.peek() != '^' {
			break
		}
		sl.next()
		// check only '^'
		if sl.isEOF() { // like "l4^"
			total += defLen
		}
	}
	return total
}

func (song *Song) NToStep(n int) int {
	if n == 0 {
		return 0
	}
	tb := song.Timebase
	step := int(float64(tb) * 4.0 / float64(n))
	return step
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
func (song *Song) CurTrack() *Track {
	for song.TrackNo >= len(song.Tracks) {
		tr := NewTrack(song.TrackNo%16, song.Timebase)
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
