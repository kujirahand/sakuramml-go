package song

import (
    "fmt"
    "sort"
    "sakuramml/event"
)

type Event struct {
    Time int
    Type event.Type
    Data1 int
    Data2 int
}

func (event *Event) GetDataBytes() []byte {
    buf := make([]byte, 3)
    buf[0] = event.Type
    buf[1] = event.Data1
    buf[2] = event.Data2
    return buf
}

type Track struct {
    Channel int
    Length int
    Octave int
    Qgate int
    Velocity int
    TimePtr int
    Events []*Event
}

func (track *Track) Init(channel int, timebase int) {
    track.Events = []*Event {}
    track.Channel = channel
    track.Length = timebase
    track.Qgate = 80
    track.Velocity = 100
}

func (track *Track) sortEvent() {
    events := track.Events
    sort.SliceStable(track.Events,
        func(i, j int) bool {
            return events[i].Time < events[j].Time
        })
}

func (track *Track) ToString() string {
    s := fmt.Sprintf("|-channel=%d", track.Channel + 1) + "\n"
    s = s + fmt.Sprintf("|-event.length=%d", len(track.Events)) + "\n"
    return s
}

type Song struct {
    Timebase int
    TrackNo int
    Tracks []*Track
}

func (song *Song) Init() {
    song.Timebase = 96
    song.Tracks = []*Track {}
    // create default track
    for i := 0; i < 16; i++ {
        track := Track{}
        track.Init(i, song.Timebase)
        song.Tracks = append(song.Tracks, &track)
    }
    song.TrackNo = 0
}

func (song *Song) ToString() string {
    s := "Timebase=" + fmt.Sprintf("%d", song.Timebase) + "\n"
    for i := 0; i < len(song.Tracks); i++ {
        track := song.Tracks[i]
        if len(track.Events) == 0 { continue }
        s += fmt.Sprintf("+Track=%d\n", (i + 1))
        s += track.ToString()
    }
    return s
}







