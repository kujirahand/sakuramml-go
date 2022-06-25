package sakuramml

import (
	"testing"
)

func TestTrackSortEvents(t *testing.T) {
	trk := NewTrack(0, 96)
	trk.AddNoteOn(96, 60, 127, 10)  // on & off
	trk.AddNoteOn(120, 30, 100, 10) // on & off
	trk.AddNoteOn(0, 30, 100, 10)   // on & off
	trk.SortEvents()
	if trk.Events[0].Time != 0 {
		t.Errorf("TrackSortEvents Error 0: %+v", trk.Events)
	}
	if trk.Events[2].Time != 96 {
		t.Errorf("TrackSortEvents Error 1: %+v", trk.Events)
	}
	if trk.Events[4].Time != 120 {
		t.Errorf("TrackSortEvents Error 2: %+v", trk.Events)
	}
}

/*
func TestAddTempoEvent(t *testing.T) {
	// tempo = 120 ... 0xFF 0x51 0x03 0x07 0xA1 0x20
	trk := NewTrack(0, 96)
	eve := trk.AddTempo(0, 120)
	actual := *eve.SongBytesToHex(eve.ExData)
	expect := "ff510307a120"
	if actual != expect {
		t.Errorf("AddTempo Error 1: %s != %s", actual, expect)
	}
	// 30 ... FF  51 03 1E 84 80
	eve = trk.AddTempo(0, 30)
	actual = BytesToHex(eve.ExData)
	expect = "ff51031e8480"
	if actual != expect {
		t.Errorf("AddTempo Error 2: %s != %s", actual, expect)
	}
	// 300 ... FF  51 03 03 0D  40
	eve = trk.AddTempo(0, 300)
	actual = BytesToHex(eve.ExData)
	expect = "ff5103030d40"
	if actual != expect {
		t.Errorf("AddTempo Error 3: %s != %s", actual, expect)
	}
}
*/

func TestAddPitchBend(t *testing.T) {
	trk := NewTrack(0, 96)
	pb := trk.AddPitchBend(0, 0)
	act := BytesToHex(pb.GetDataBytes())
	exp := "e00040"
	if act != exp {
		t.Errorf("AddPitchBend Error p%%0 : %s != %s", act, exp)
	}

	pb = trk.AddPitchBend(0, 100)
	act = BytesToHex(pb.GetDataBytes())
	exp = "e06440"
	if act != exp {
		t.Errorf("AddPitchBend Error p%%100 : %s != %s", act, exp)
	}
}
