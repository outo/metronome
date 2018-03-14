package metronome

import "github.com/outo/metronome/param"

//the intended purpose is to delegate to performer at a rate described by bpm
type Metronome func(bpm param.Bpm, performer BeatPerformer)

//provides a volume of ambient environment
type VolumeMeter func() (volume int)


//will "perform" beat
// count parameter can be used to decide if Tick or Tock should be displayed
// volume decides how strong the beat is going to be
type BeatPerformer func(beatCount, volume int)
