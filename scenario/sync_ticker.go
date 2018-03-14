package scenario

import (
	"github.com/outo/metronome/generator"
	"github.com/outo/metronome/track"
	"github.com/outo/metronome/param"
	"github.com/outo/metronome/chart"
	"time"
	"github.com/outo/metronome"
)

//similar implementation to adaptive delay but with the use of Ticker, won't fire immediately
func synchronousTicker(volumeMeter metronome.VolumeMeter, numberOfBeats int, _ track.Tracker) metronome.Metronome {
	return func(bpm param.Bpm, performer metronome.BeatPerformer) {

		beatCount := 0

		//create Ticker instance, initialised with desired interval
		ticker := time.NewTicker(bpm.Interval())

		//defer ensures the Ticker stops sending messages before we leave this function
		defer ticker.Stop()

		//for each Ticker's timing message
		for range ticker.C {

			volume := volumeMeter()

			performer(beatCount, volume)

			if beatCount++; beatCount >= numberOfBeats {
				break
			}
		}
	}
}

func SynchronousTicker(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
	demonstrateScenario("sync_ticker", generator.ShortDuration, bpm, numberOfBeats, false, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			synchronousTicker(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Synchronous metronome with <b>ticker</b>. " +
					"Notice initial delay before first volume measurement, caused by Ticker <b>not firing immediately</b>. ",
				Groups:      defaultGroups,
			}
		},
	)
}
