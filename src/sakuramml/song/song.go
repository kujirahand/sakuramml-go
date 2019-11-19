package song

import (
	"fmt"
	"sakuramml/event"
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
func (track *Track) AddNoteOn(time, note, vel, lenStep int) {
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
}

// AddCC Add Controll Change
func (track *Track) AddCC(time, no, value int) {
	cc := event.Event{
		Time:      time,
		ByteCount: 3,
		Type:      event.CC | track.Channel,
		Data1:     no,
		Data2:     value,
	}
	track.AddEvent(cc)
}

// AddProgramChange Add Controll Change
func (track *Track) AddProgramChange(time, value int) {
	pc := event.Event{
		Time:      time,
		ByteCount: 2,
		Type:      event.ProgramChange | track.Channel,
		Data1:     value,
	}
	track.AddEvent(pc)
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

// Song is info of song, include tracks
type Song struct {
	Timebase int
	TrackNo  int
	Stack    []interface{}
	Tracks   []*Track
	Debug    bool
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
	s.TrackNo = 0
	s.Stack = make([]interface{}, 0, 256)
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
