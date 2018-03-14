package scenario

import (
	"github.com/outo/metronome/track"
	"github.com/outo/metronome/param"
	"github.com/outo/metronome/chart"
	"time"
	"github.com/outo/metronome/generator"
	"github.com/outo/metronome"
)

// order of beats can be distorted
// volume measurement may happen simultaneously (can sampler cope?)
// not all beats have been generated before program ended
// go with literal function closure - danger!
func asynchronousButNaive(volumeMeter metronome.VolumeMeter, numberOfBeats int, _ track.Tracker) metronome.Metronome {
	return func(bpm param.Bpm, performer metronome.BeatPerformer) {

		ticker := time.NewTicker(bpm.Interval())
		defer ticker.Stop()

		for beatCount := 0; beatCount < numberOfBeats; beatCount ++ {

			go func(beatCount int) {
				volume := volumeMeter()

				performer(beatCount, volume)
			}(beatCount)

			<-ticker.C
		}
	}
}

func AsynchronousButNaive(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
	demonstrateScenario("async_naive", generator.RandomisedWithProportionalCapsDuration, bpm, numberOfBeats, true, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			asynchronousButNaive(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Asynchronous metronome - naive implementation. " +
					"Execution time of volume measurement is <b>random</b>. ",
				Groups: map[string]chart.Group{
					"ghost":  ghostGroup,
					"volume": volumeGroup.WithContent("Volume measurement <br> (notice the overlaps)"),
					"beat":   beatGroup,
				},
			}
		},
	)
}
