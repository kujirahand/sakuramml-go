package midi

import (
    "testing"
    "sakuramml/midi"
    "sakuramml/utils"
)

func TestGetDeltaTimeBytes(t *testing.T) {
    // 10
    a := utils.BytesToHex(midi.GetDeltaTimeBytes(10))
    e := utils.BytesToHex([]byte { 10 })
    if a != e {
        t.Errorf("GetDeltaTime:%s != %s", a, e)
    }
    // 0x81
    a2 := utils.BytesToHex(midi.GetDeltaTimeBytes(0x81))
    e2 := utils.BytesToHex([]byte { 0x81, 01 })
    if a2 != e2 {
        t.Errorf("GetDeltaTime:%s != %s", a2, e2)
    }
    // 0xFFFFFFF
    a3 := utils.BytesToHex(midi.GetDeltaTimeBytes(0x1FFFFF))
    e3 := utils.BytesToHex([]byte { 0xFF, 0xFF, 0x7F })
    if a3 != e3 {
        t.Errorf("GetDeltaTime:%s != %s", a3, e3)
    }
    // 0xFFFFFFF
    a4 := utils.BytesToHex(midi.GetDeltaTimeBytes(0x8000000))
    e4 := utils.BytesToHex([]byte { 0xC0, 0x80, 0x80, 0x00 })
    if a4 != e4 {
        t.Errorf("GetDeltaTime:%s != %s", a4, e4)
    }
}

