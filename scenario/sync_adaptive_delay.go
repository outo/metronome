package scenario

import (
	"github.com/outo/metronome/generator"
	"github.com/outo/metronome/track"
	"github.com/outo/metronome/param"
	"github.com/outo/metronome/chart"
	"time"
	"github.com/outo/metronome"
)

//adaptive delay mitigates unwanted delay as long as the measurement does not take anywhere near interval duration
func synchronousAdaptiveDelay(volumeMeter metronome.VolumeMeter, numberOfBeats int, _ track.Tracker) metronome.Metronome {
	return func(bpm param.Bpm, performer metronome.BeatPerformer) {

		for beatCount := 0; beatCount < numberOfBeats; beatCount ++ {

			//let's find out the duration of volume measurement...
			start := time.Now()

			volume := volumeMeter()

			performer(beatCount, volume)

			//...so that we can adapt the planned delay accordingly (or skip it if took too long)
			if adaptiveDelay := bpm.Interval() - time.Since(start); adaptiveDelay > 0 {
				time.Sleep(adaptiveDelay)
			}
		}
	}
}

func SynchronousAdaptiveDelay(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
	demonstrateScenario("sync_adaptive_delay", generator.ShortDuration, bpm, numberOfBeats, false, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			synchronousAdaptiveDelay(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Synchronous metronome with <b>adaptive</b> delay and volume measurement before each beat. ",
				Groups:      defaultGroups,
			}
		},
	)

	demonstrateScenario("sync_adaptive_delay_oscillating_execution_time", generator.TriangularPatternDuration, bpm, numberOfBeats, false, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			synchronousAdaptiveDelay(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Synchronous metronome with <b>adaptive</b> delay and volume measurement before each beat. " +
					"Execution time of volume measurement is <b>oscillating</b>. ",
				Groups:   defaultGroups,
				CodeBase: "sync_adaptive_delay",
			}
		},
	)
}
