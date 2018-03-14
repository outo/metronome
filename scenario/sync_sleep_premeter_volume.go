package scenario

import (
	"github.com/outo/metronome/generator"
	"github.com/outo/metronome/track"
	"github.com/outo/metronome/param"
	"time"
	"github.com/outo/metronome"
	"github.com/outo/metronome/chart"
)

//fulfills well simple requirement with single measurement before beats performing
func synchronousSleepWithPreMeter(volumeMeter metronome.VolumeMeter, numberOfBeats int, _ track.Tracker) metronome.Metronome {
	return func(bpm param.Bpm, performer metronome.BeatPerformer) {

		//measure the volume before the beats are to be performed
		volume := volumeMeter()

		//loop over beatCount values ranging [0, numberOfBeats)
		for beatCount := 0; beatCount < numberOfBeats; beatCount ++ {

			//delegation to beat performer which will (as an example) print Ticks and Tocks
			performer(beatCount, volume)

			//planned delay so that the beats appear at equal intervals
			//Bpm.Interval() will perform simple calculation of interval length
			// to match required bpm value.
			time.Sleep(bpm.Interval())
		}
	}
}

func SynchronousSleepWithPreMeter(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
	demonstrateScenario("sync_sleep_premeter_volume", generator.ShortDuration, bpm, numberOfBeats, false, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			synchronousSleepWithPreMeter(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Synchronous metronome with <b>constant</b> delay and <b>single volume measurement</b>.",
				Groups:      defaultGroups,
			}
		},
	)
}
