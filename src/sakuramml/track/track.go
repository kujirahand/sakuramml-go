package track

import (
	"fmt"
	"sakuramml/event"
	"sakuramml/utils"
	"sort"
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
	PitchBend int
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
	track.PitchBend = 0
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

// AddPitchBendEasy func ... p command / 簡易ピッチベンドを書き込む(0~63~127の範囲)
func (track *Track) AddPitchBendEasy(time, value int) *event.Event {
	v := value*128 - 8192
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
