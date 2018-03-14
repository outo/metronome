package scenario

import (
   "sync"
   "github.com/outo/metronome/generator"
   "github.com/outo/metronome/track"
   "github.com/outo/metronome/param"
   "github.com/outo/metronome/chart"
   "time"
   "github.com/outo/metronome"
)

//because this struct and demo code resides in the same package it may be tempting
// to access struct's field directly - don't! unless you synchronise access to them
type SharedState struct {
   mutex                sync.Mutex
   volume               int
   firstMeasurementDone bool
   terminationRequested bool
   measuringFinished    bool
}

func NewSharedState() *SharedState {
   return &SharedState{
      volume: -1,
   }
}

// convenience method reducing repetition around locking/unlocking from 2 to 1 line
func (s *SharedState) lock() *SharedState {
   s.mutex.Lock()
   return s
}

// convenience method reducing repetition around locking/unlocking from 2 to 1 line
func (s *SharedState) unlock() {
   s.mutex.Unlock()
}

func (s *SharedState) KeepMeasuring() bool {
   defer s.lock().unlock() // .lock() invoked straightaway, .Unlock() will defer
   return !s.terminationRequested
}

func (s *SharedState) RequestTermination() {
   defer s.lock().unlock()
   s.terminationRequested = true
}

func (s *SharedState) VolumeMeasuredAtLeastOnce() bool {
   defer s.lock().unlock()
   return s.volume != -1
}

func (s *SharedState) NewVolumeMeasurement(volume int) {
   defer s.lock().unlock()
   s.volume = volume
}

func (s *SharedState) MeasuringFinished() {
   defer s.lock().unlock()
   s.measuringFinished = true
}

func (s *SharedState) HasMeasuringFinished() bool {
   defer s.lock().unlock()
   return s.measuringFinished
}

func (s *SharedState) LatestVolume() int {
   defer s.lock().unlock()
   return s.volume
}

func asynchronousWithMutex(volumeMeter metronome.VolumeMeter, numberOfBeats int, _ track.Tracker) metronome.Metronome {
   return func(bpm param.Bpm, performer metronome.BeatPerformer) {

      sharedState := NewSharedState()

      // continually and sequentially measure volume
      go func() {
         for sharedState.KeepMeasuring() {
            sharedState.NewVolumeMeasurement(volumeMeter())
         }
         sharedState.MeasuringFinished()
      }()

      // wait for first measurement to come through
      for !sharedState.VolumeMeasuredAtLeastOnce() {
      }

      ticker := time.NewTicker(bpm.Interval())
      defer ticker.Stop()

      for beatCount := 0; beatCount < numberOfBeats; beatCount ++ {
         performer(beatCount, sharedState.LatestVolume())
         <-ticker.C

      }

      sharedState.RequestTermination()

      for !sharedState.HasMeasuringFinished() {
      }

   }
}

func AsynchronousWithMutex(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
   demonstrateScenario("async_with_mutex", generator.RandomisedWithProportionalCapsDuration, bpm, numberOfBeats, false, metas, eventsChannel,
      func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
         asynchronousWithMutex(volumeMeter, numberOfBeats, tracker)(bpm, performer)
         return chart.Meta{
            Description: "Asynchronous metronome with <b>mutex</b> protecting shared state. " +
               "Execution time of volume measurement is <b>random</b>. ",
            Groups: defaultGroups,
         }
      },
   )
}
