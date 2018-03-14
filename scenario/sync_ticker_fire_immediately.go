package scenario

import (
	"github.com/outo/metronome/generator"
	"github.com/outo/metronome/track"
	"github.com/outo/metronome/param"
	"github.com/outo/metronome/chart"
	"time"
	"github.com/outo/metronome"
)

//implementation with ticker, fires immediately
func synchronousTickerFireImmediately(volumeMeter metronome.VolumeMeter, numberOfBeats int, _ track.Tracker) metronome.Metronome {
	return func(bpm param.Bpm, performer metronome.BeatPerformer) {

		ticker := time.NewTicker(bpm.Interval())
		defer ticker.Stop()

		for beatCount := 0; beatCount < numberOfBeats; beatCount ++ {

			volume := volumeMeter()

			//my preference is to wait for the right time here, with pre-fetched volume
			// because timing feels more important in our case than freshness of volume
			<-ticker.C

			performer(beatCount, volume)

			//<-ticker.C //alternative position
		}
	}
}

func SynchronousTickerFireImmediately(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
	demonstrateScenario("sync_ticker_fire_immediately", generator.ShortDuration, bpm, numberOfBeats, false, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			synchronousTickerFireImmediately(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Synchronous metronome with <b>ticker firing immediately</b>. ",
				Groups:      defaultGroups,
			}
		},
	)
	demonstrateScenario("sync_ticker_fire_immediately_oscillating_execution_time", generator.TriangularPatternDuration, bpm, numberOfBeats, false, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			synchronousTickerFireImmediately(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Synchronous metronome with <b>ticker firing immediately</b>. " +
					"Execution time of volume measurement is <b>oscillating</b>. ",
				Groups:   defaultGroups,
				CodeBase: "sync_ticker_fire_immediately",
			}
		},
	)
	demonstrateScenario("sync_ticker_fire_immediately_randomised_execution_time", generator.RandomisedWithProportionalCapsDuration, bpm, numberOfBeats, false, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			synchronousTickerFireImmediately(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Synchronous metronome with <b>ticker firing immediately</b>. " +
					"Execution time of volume measurement is <b>random</b>. ",
				Groups:   defaultGroups,
				CodeBase: "sync_ticker_fire_immediately",
			}
		},
	)
}
