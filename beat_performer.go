package metronome

import (
	"github.com/outo/metronome/track"
	"fmt"
	"strconv"
)

//this beat performer tracks the beats (not actually performing)
// normally it would play or print beats but here we care about tracking only
// so that we can plot a timing chart
func CreateBeatPerformerWithTrackingWithIOContent(tracker track.Tracker) BeatPerformer {
	var beats = []string{"I", "O"}

	return func(count, power int) {
		tickOrTock := beats[count%2]

		event := track.Event{
			Category:    "beat",
			Description: fmt.Sprintf("beat #%d at %d power", count, power),
			Content:     tickOrTock,
		}
		tracker(event)
	}
}

func CreateBeatPerformerWithTrackingWithCount(tracker track.Tracker) BeatPerformer {
	return func(count, power int) {
		event := track.Event{
			Category:    "beat",
			Description: fmt.Sprintf("beat #%d at %d power", count, power),
			Content:     strconv.Itoa(count),
		}
		tracker(event)
	}
}

