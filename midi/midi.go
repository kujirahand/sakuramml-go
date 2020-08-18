package midi

import (
	// "fmt"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"github.com/kujirahand/sakuramml-go/song"
	"github.com/kujirahand/sakuramml-go/track"
)

// GetUint16 make 2 bytes BigEndian data
func GetUint16(v int) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(v))
	return buf
}

// GetUint32 make 4 bytes BigEndian data
func GetUint32(v int) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v))
	return buf
}

// Save song to midi stream
func Save(s *song.Song, w io.Writer) {
	midiformat := 1

	// count track count
	var trackCount = 0
	for _, track := range s.Tracks {
		if len(track.Events) == 0 {
			continue
		}
		trackCount++
	}

	// Write header
	w.Write([]byte("MThd"))
	w.Write(GetUint32(6))
	w.Write(GetUint16(midiformat))
	w.Write(GetUint16(trackCount))
	w.Write(GetUint16(s.Timebase))

	// Write tracks
	for _, track := range s.Tracks {
		if len(track.Events) == 0 {
			continue
		}
		block := getTrackData(track)
		// track header
		w.Write([]byte("MTrk"))
		w.Write(GetUint32(len(block)))
		w.Write(block)
		// fmt.Println(i, track)
	}
}

// SaveToFile : song save to midi file
func SaveToFile(sng *song.Song, outfile string) {
	f, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	Save(sng, f)
}

func getTrackData(track *track.Track) []byte {
	buf := new(bytes.Buffer)
	track.SortEvents()
	events := track.Events
	pTime := 0
	for _, event := range events {
		dtime := event.Time - pTime
		if dtime < 0 {
			dtime = 0
		}
		buf.Write(GetDeltaTimeBytes(dtime))
		buf.Write(event.GetDataBytes())
		pTime = event.Time
	}
	// EOT(End Of Track)のイベントを書く
	buf.Write([]byte{0, 0xFF, 0x2F, 00})
	return buf.Bytes()
}

// GetDeltaTimeBytes : make Delta time
func GetDeltaTimeBytes(v int) []byte {
	if v == 0 {
		b := make([]byte, 1)
		return b
	}
	var buf [256]byte
	i := 0
	for v > 0 {
		if i > 255 {
			log.Fatal("Time value overflow")
		}
		buf[i] = byte(v & 0x7F)
		if i > 0 {
			buf[i] = 0x80 | buf[i]
		}
		i++
		v = v >> 7
	}
	var out = make([]byte, i)
	for j := 0; j < i; j++ {
		out[i-j-1] = buf[j]
	}
	return out
}
