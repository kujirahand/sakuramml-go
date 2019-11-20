package song

import (
	"fmt"
	"sakuramml/event"
	"sakuramml/utils"
	"sort"
	"strconv"
)

const (
	// QgateModeStep for QGateMode
	QgateModeStep = "step"
	// QgateModeRate for QGateMode
	QgateModeRate = "rate"
)

// Track is info of track
type Track struct {
	Channel   int
	Length    int // step
	Octave    int
	Qgate     int    // ref: QgateMode
	QgateMode string // step or rate
	Velocity  int
	Time      int
	Events    []event.Event
}

// NewTrack func
func NewTrack(channel int, timebase int) *Track {
	track := Track{}
	track.Events = make([]event.Event, 0, 256) // Default Event
	track.Channel = channel
	track.Length = timebase
	track.Qgate = 80
	track.QgateMode = QgateModeRate
	track.Velocity = 100
	track.Octave = 5
	return &track
}

// AddEvent func
func (track *Track) AddEvent(event event.Event) {
	track.Events = append(track.Events, event)
}

// AddNoteOn add NoteOn event to track
func (track *Track) AddNoteOn(time, note, vel, lenStep int) (*event.Event, *event.Event) {
	eon := event.Event{
		Time:      time,
		ByteCount: 3,
		Type:      event.NoteOn | track.Channel,
		Data1:     note,
		Data2:     vel,
	}
	eoff := event.Event{
		Time:      time + lenStep,
		ByteCount: 3,
		Type:      event.NoteOff | track.Channel,
		Data1:     note,
		Data2:     vel,
	}
	track.AddEvent(eon)
	track.AddEvent(eoff)
	return &eon, &eoff
}

// AddCC Add Controll Change
func (track *Track) AddCC(time, no, value int) *event.Event {
	cc := event.Event{
		Time:      time,
		ByteCount: 3,
		Type:      event.CC | track.Channel,
		Data1:     no,
		Data2:     value,
	}
	track.AddEvent(cc)
	return &cc
}

// AddProgramChange Add Controll Change
func (track *Track) AddProgramChange(time, value int) *event.Event {
	pc := event.Event{
		Time:      time,
		ByteCount: 2,
		Type:      event.ProgramChange | track.Channel,
		Data1:     value,
	}
	track.AddEvent(pc)
	return &pc
}

// AddPitchBend func ... p% command (%をつけると-8192~0~8191))
func (track *Track) AddPitchBend(time, value int) *event.Event {
	// calc msb, lsb
	v := value + 8192
	v = utils.InRange(0, v, 16383)
	msb := v >> 7 & 0x7f
	lsb := v & 0x7f
	// gen
	pb := event.Event{
		Time:      time,
		ByteCount: 3,
		Type:      event.PitchBend | track.Channel,
		Data1:     lsb, // low byte <--- MIDI仕様からすると逆に思えるが lsb, msb の順が正しい
		Data2:     msb, // high byte
	}
	track.AddEvent(pb)
	return &pb
}

// AddPitchBendEx func ... p command / 簡易ピッチベンドを書き込む(0~63~127の範囲)
func (track *Track) AddPitchBendEx(time, value int) *event.Event {
	v := value * 64
	return track.AddPitchBend(time, v)
}

// AddTempo func
func (track *Track) AddTempo(time, tempo int) *event.Event {
	e := event.Event{
		Time:      time,
		ByteCount: 6,
		Type:      event.Tempo,
	}
	mpq := uint32(60000000 / tempo)
	e.ExData = []byte{
		0xFF,
		0x51,
		0x03,
		byte(mpq >> 16 & 0xff),
		byte(mpq >> 8 & 0xff),
		byte(mpq & 0xff),
	}
	track.AddEvent(e)
	return &e
}

// AddMeta func
func (track *Track) AddMeta(time, metaType int, data []byte) *event.Event {
	e := event.Event{
		Time:      time,
		ByteCount: 3 + len(data),
		Type:      0xFF00 | metaType,
	}
	buf := make([]byte, 3+len(data))
	e.ExData = buf
	buf[0] = 0xFF
	buf[1] = byte(metaType)
	buf[2] = byte(len(data))
	for i := 0; i < len(data); i++ {
		buf[3+i] = data[i]
	}
	track.AddEvent(e)
	return &e
}

// SortEvents sort Events of track
func (track *Track) SortEvents() {
	events := track.Events
	sort.SliceStable(track.Events,
		func(i, j int) bool {
			return events[i].Time < events[j].Time
		})
}

// ToString convert to string
func (track *Track) ToString() string {
	s := fmt.Sprintf("|-channel=%d", track.Channel+1) + "\n"
	s = s + fmt.Sprintf("|-event.length=%d", len(track.Events)) + "\n"
	return s
}

// LoopItem struct
type LoopItem struct {
	Index     int
	Count     int
	BeginNode interface{} // Node
	EndNode	  interface{} // Node
}

// Song is info of song, include tracks
type Song struct {
	Timebase  int
	Tempo     int
	TrackNo   int
	Stack     []interface{} // values stack
	Tracks    []*Track
	Debug     bool
	LoopStack []*LoopItem
	MoveNode  interface{} // Node
}

// NewSong func
func NewSong() *Song {
	s := Song{}
	s.Debug = false
	s.Timebase = 96
	s.Tracks = []*Track{}
	// create default track
	for i := 0; i < 16; i++ {
		track := NewTrack(i, s.Timebase)
		s.Tracks = append(s.Tracks, track)
	}
	s.Tempo = 120 // default Tempo (but not write)
	s.TrackNo = 0
	s.Stack = make([]interface{}, 0, 256)
	s.MoveNode = nil
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
		return song.LoopStack[ilen - 1]
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
