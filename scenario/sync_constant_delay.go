package scenario

import (
	"github.com/outo/metronome/generator"
	"github.com/outo/metronome/track"
	"github.com/outo/metronome/param"
	"time"
	"github.com/outo/metronome"
	"github.com/outo/metronome/chart"
)

//when requirement changes to measure the volume in real-time it introduces delays
func synchronousConstantDelay(volumeMeter metronome.VolumeMeter, numberOfBeats int, _ track.Tracker) metronome.Metronome {
	return func(bpm param.Bpm, performer metronome.BeatPerformer) {

		for beatCount := 0; beatCount < numberOfBeats; beatCount ++ {
			volume := volumeMeter()

			performer(beatCount, volume)

			time.Sleep(bpm.Interval())
		}
	}
}

func SynchronousConstantDelay(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
	demonstrateScenario("sync_constant_delay", generator.ShortDuration, bpm, numberOfBeats, false, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			synchronousConstantDelay(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Synchronous metronome with <b>constant</b> delay and volume measurement before each beat.",
				Groups:      defaultGroups,
			}
		},
	)
}
