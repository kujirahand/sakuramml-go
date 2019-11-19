package utils

import (
	"fmt"
)

// BytesToHex func
func BytesToHex(b []byte) string {
	s := ""
	for _, v := range b {
		s += fmt.Sprintf("%02x", v)
	}
	return s
}

// MidiRange func
func MidiRange(v int) int {
	if v < 0 {
		v = 0
	} else if v > 127 {
		v = 127
	}
	return v
}
