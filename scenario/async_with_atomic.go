package scenario

import (
	"github.com/outo/metronome/generator"
	"github.com/outo/metronome/track"
	"github.com/outo/metronome/param"
	"github.com/outo/metronome/chart"
	"sync/atomic"
	"time"
	"github.com/outo/metronome"
)

// these constants help us determine the stage our program is in,
// variable sharedStage will be set to one of them
const (
	// initially, we haven't received our first volume measurement, so we need to wait
	StageNeedVolume      int32 = iota

	// once we measure volume for the first time, this will indicate we can progress to beat performing
	StageFirstVolumeProvided

	// indicate that we have performed enough beats and would like to finish gracefully
	StageNeedTermination

	// indicate that the goroutine measuring volume has terminated
	StageTerminated
)

func createRealTimeVolumeVariable(
	sharedStage *int32,
	volumeMeter metronome.VolumeMeter) (*int32) {

	sharedVolumeReading := int32(-1)

	// continually and sequentially measure volume
	go func() {
		// safely read the sharedStage value and execute the for-loop body, until we need to terminate
		for atomic.LoadInt32(sharedStage) != StageNeedTermination {

			// safe volume is locally-scoped variable (not shared), which will temporarily hold volume value
			safeVolume := int32(volumeMeter())

			// atomically set the sharedVolumeReading to safeVolume and
			// store old value of sharedVolumeReading in safeOldVolume
			// (sharedVolumeReading is set to -1 before this loop)
			safeOldVolume := atomic.SwapInt32(&sharedVolumeReading, safeVolume)

			// on the first iteration of this loop the following will evaluate to true
			if safeOldVolume == -1 {
				// indicate we have measure volume for the first time, go to next stage
				atomic.StoreInt32(sharedStage, StageFirstVolumeProvided)
			}
		}

		// indicate this Goroutine is just about to be done
		atomic.StoreInt32(sharedStage, StageTerminated)
	}()

	return &sharedVolumeReading
}

func asynchronousWithAtomic(
	volumeMeter metronome.VolumeMeter,
	numberOfBeats int,
	_ track.Tracker) metronome.Metronome {

	return func(bpm param.Bpm, performer metronome.BeatPerformer) {

		// initial stage - we need first volume measurement before we start beating
		sharedStage := StageNeedVolume

		// this call will create a Goroutine that will continually provide the latest volume in a
		// variable pointed to by sharedVolume
		sharedVolume := createRealTimeVolumeVariable(&sharedStage, volumeMeter)

		// this loop will wait for when the stage indicates we have first volume measurement
		for atomic.LoadInt32(&sharedStage) != StageFirstVolumeProvided {
		}

		// standard Ticker
		ticker := time.NewTicker(bpm.Interval())
		// and deferred tear down of it
		defer ticker.Stop()

		for beatCount := 0; beatCount < numberOfBeats; beatCount ++ {

			// safely fetch the value of variable referenced by sharedVolume
			// it will not block but retrieve whatever was the last value
			safeVolume := int(atomic.LoadInt32(sharedVolume))

			// standard beat performance
			performer(beatCount, safeVolume)

			// wait for next timing message
			<-ticker.C
		}

		// indicate that we have performed enough beats and would like
		// to terminate volume measurement. That will allow us to gracefully and
		// in controlled manner terminate program.
		atomic.StoreInt32(&sharedStage, StageNeedTermination)

		// wait until volume measurement Goroutine is done
		for atomic.LoadInt32(&sharedStage) != StageTerminated {
		}

	}
}

func AsynchronousWithAtomic(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
	demonstrateScenario("async_with_atomic", generator.RandomisedWithProportionalCapsDuration, bpm, numberOfBeats, false, metas, eventsChannel,
		func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
			asynchronousWithAtomic(volumeMeter, numberOfBeats, tracker)(bpm, performer)
			return chart.Meta{
				Description: "Asynchronous metronome with <b>atomic variables</b>. " +
					"Execution time of volume measurement is <b>random</b>. ",
				Groups: defaultGroups,
			}
		},
	)
}
