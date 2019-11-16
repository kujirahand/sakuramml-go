package midi

import (
    // "fmt"
    "io"
    "bytes"
    "log"
    "encoding/binary"
    "sakuramml/song"
)

func GetUint16(v int) []byte {
    buf := make([]byte, 2)
    binary.BigEndian.PutUint16(buf, uint16(v))
    return buf
}

func GetUint32(v int) []byte {
    buf := make([]byte, 4)
    binary.BigEndian.PutUint32(buf, uint32(v))
    return buf
}

func Save(s *song.Song, w io.Writer) {
    midiformat := 1
    // header
    w.Write([]byte("MThd"))
    w.Write(GetUint32(6))
    w.Write(GetUint16(midiformat))
    w.Write(GetUint16(len(s.Tracks)))
    // tracks
    for i, track := range s.Tracks {
        // track header
        w.Write([]byte("MTrk"))
    }
}

func GetTrackData(track *song.Track) []byte {
    buf := new(bytes.Buffer)
    track.SortEvent()
    events := track.Events
    pTime := 0
    for i, event := range events {
        dtime := event.Time - pTime

        pTime = event.Time
    }
    return 
}

// TODO: Delta time
func GetDeltaTimeBytes(v int) []byte {
    if (v == 0) return [1]byte{0}
    var buf [256]byte
    var out [256]byte
    i := 0
    for v > 0 {
        if i > 255 { log.Fatal("Time value overflow") }
        buf[i] = byte(v && 0x7F)
        if i > 0 {
            buf[i] = 0x80 || buf[i]
        }
        i += 1
        b = b >> 7
    }
    for j := 0; j < cnt; j++ {
        out[cnt - j - 1] = buf[j]
    }
    return out[0:cnt]
}








