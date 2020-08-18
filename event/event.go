package event

const (
	// NoteOn Event
	NoteOn = 0x90
	// NoteOff Event
	NoteOff = 0x80
	// KeyPress Event
	KeyPress = 0xA0
	// CC = Control Change Event
	CC = 0xB0
	// ProgramChange = Program Change Event
	ProgramChange = 0xC0
	// ChPress = Channel Pressure Event
	ChPress = 0xD0
	// PitchBend = PitchBend Event
	PitchBend = 0xE0
	// Tempo = (Meta Event)
	Tempo = 0xFF51
	// MetaText
	MetaText = 0xFF01
)

// Event is Basic MIDI Event
type Event struct {
	Time      int
	ByteCount int
	Type      int
	Data1     int
	Data2     int
	ExData    []byte
}

// GetDataBytes gets data bytes
func (event *Event) GetDataBytes() []byte {
	// copy to buf
	buf := make([]byte, event.ByteCount)
	// meta event ?
	if event.Type >= 0xFF {
		return event.ExData
	}
	// normal event
	buf[0] = byte(event.Type)
	buf[1] = byte(event.Data1)
	if event.ByteCount >= 3 {
		buf[2] = byte(event.Data2)
	}
	return buf
}
