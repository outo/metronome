package scenario

import (
	"github.com/outo/metronome/generator"
	"github.com/outo/metronome/track"
	"github.com/outo/metronome/param"
	"github.com/outo/metronome/chart"
	"sync"
	"time"
	"github.com/outo/metronome"
)

func asynchronousStillNaiveWithWaitGroup(volumeMeter metronome.VolumeMeter, numberOfBeats int, _ track.Tracker) metronome.Metronome {
	return func(bpm param.Bpm, performer metronome.BeatPerformer) {

		taskGroup := sync.WaitGroup{}
		defer taskGroup.Wait()

		ticker := time.NewTicker(bpm.Interval())
		defer ticker.Stop()

		for beatCount := 0; beatCount < numberOfBeats; beatCount ++ {

			taskGroup.Add(1)
			go func(beatCount int) {
				volume := volumeMeter()

				performer(beatCount, volume)
				taskGroup.Done()
			}(beatCount)

			<-ticker.C
		}

	}
}

func AsynchronousStillNaiveWithWaitGroup(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
	demonstrateScenario("async_naive_waiting_for_completion", generator.RandomisedWithProportionalCapsDuration, bpm, numberOfBeats, true, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			asynchronousStillNaiveWithWaitGroup(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Asynchronous metronome - naive implementation that at least <b>waits for the completion</b> of the whole execution. "+
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
