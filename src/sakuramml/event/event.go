package event

type Type int

const (
    NoteOn    = 0x90
    NoteOff   = 0x80
    KeyPress  = 0xA0
    CC        = 0xB0
    Program   = 0xC0
    ChPress   = 0xD0
    PitchBend = 0xE0
)


