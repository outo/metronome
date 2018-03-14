package chart

import (
	"github.com/outo/metronome/param"
	"github.com/outo/metronome/track"
	"fmt"
	"time"
)

func ApplyGhosting(bpm param.Bpm, numberOfBeats int, tracker track.Tracker) {
	var beats = []string{"I", "O"}
	for count := 0; count < numberOfBeats; count++ {
		tickOrTock := beats[count%2]

		tracker(track.Event{
			Category:       "ghost",
			Description:    fmt.Sprintf("ghost #%d", count),
			Content:        tickOrTock,
			TimeSinceStart: time.Duration(count) * bpm.Interval(),
		})
	}
}

