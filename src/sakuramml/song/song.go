package song

import (
	"fmt"
	"sort"
)

// Event is Basic MIDI Event
type Event struct {
	Time  int
	Type  int
	Data1 int
	Data2 int
}

// GetDataBytes gets data bytes
func (event *Event) GetDataBytes() []byte {
	buf := make([]byte, 3)
	buf[0] = byte(event.Type)
	buf[1] = byte(event.Data1)
	buf[2] = byte(event.Data2)
	return buf
}

// Track is info of track
type Track struct {
	Channel  int
	Length   int
	Octave   int // step
	Qgate    int // step
	Velocity int
	Time     int
	Events   []Event
}

// NewTrack func
func NewTrack(channel int, timebase int) *Track {
	track := Track{}
	track.Events = make([]Event, 0, 256) // Default Event
	track.Channel = channel
	track.Length = timebase
	track.Qgate = int(float64(timebase) * 0.8)
	track.Velocity = 100
	track.Octave = 5
	return &track
}

// AddEvent func
func (track *Track) AddEvent(event Event) {
	track.Events = append(track.Events, event)
}

// AddNoteOn add NoteOn event to track
func (track *Track) AddNoteOn(time, note, vel, lenStep int) {
	eon := Event{
		Time:  time,
		Type:  0x90 | track.Channel,
		Data1: note,
		Data2: vel,
	}
	eoff := Event{
		Time:  time + lenStep,
		Type:  0x80 | track.Channel,
		Data1: note,
		Data2: vel,
	}
	track.AddEvent(eon)
	track.AddEvent(eoff)
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
	IValues  []int
	Tracks   []*Track
}

// NewSong func
func NewSong() *Song {
	s := Song{}
	s.Timebase = 96
	s.Tracks = []*Track{}
	// create default track
	for i := 0; i < 16; i++ {
		track := NewTrack(i, s.Timebase)
		s.Tracks = append(s.Tracks, track)
	}
	s.TrackNo = 0
	s.IValues = make([]int, 0, 256)
	return &s
}

// PushIValue func
func (song *Song) PushIValue(v int) {
	song.IValues = append(song.IValues, v)
}

// PopIValue func
func (song *Song) PopIValue() int {
	ilen := len(song.IValues)
	if ilen > 0 {
		v := song.IValues[ilen-1]
		song.IValues = song.IValues[0 : ilen-1]
		return v
	}
	return 0
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
