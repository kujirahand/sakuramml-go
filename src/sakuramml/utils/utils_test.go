package utils

import (
	"testing"
)

func TestBytesToHex(t *testing.T) {
	// - 1 -
	a := BytesToHex([]byte{1, 2, 3, 4})
	e := "01020304"
	if a != e {
		t.Errorf("TestBytesToHex: %s != %s", a, e)
	}
	// - 2 -
	a2 := BytesToHex([]byte{0xFF, 2, 3, 4})
	e2 := "ff020304"
	if a2 != e2 {
		t.Errorf("TestBytesToHex: %s != %s", a2, e2)
	}
}