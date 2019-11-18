package song

import (
	"testing"
)

func TestTrackSortEvents(t *testing.T) {
	trk := Track{}
	trk.Init(0, 96)
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
