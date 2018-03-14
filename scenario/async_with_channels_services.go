package scenario

import (
   "github.com/outo/metronome/generator"
   "github.com/outo/metronome/track"
   "github.com/outo/metronome/param"
   "github.com/outo/metronome/chart"
   "github.com/outo/metronome"
   "time"
   "fmt"
)

func countingService(maxCount int) (countChannel <-chan int) {
   channel := make(chan int)
   go func() {
      for count := 0; count < maxCount; count ++ {
         channel <- count
      }
      close(channel)
   }()
   return channel
}

type timing struct{}
type terminate struct{}

func isTerminating(terminationRequestChannel <-chan terminate) bool {
   select {
   case _, ok := <-terminationRequestChannel:
      return !ok
   default:
   }
   return false
}

func timingService(
   interval time.Duration,
   startTimingChannel <-chan ready,
   terminationRequestChannel <-chan terminate) (timingChannel <-chan timing) {

   channel := make(chan timing)
   go func() {

      <-startTimingChannel

      channel <- timing{}

      ticker := time.NewTicker(interval)

      for !isTerminating(terminationRequestChannel) {
         <-ticker.C
         channel <- timing{}
      }

      fmt.Println("Hurray!")
      ticker.Stop()
      close(channel)
   }()

   return channel
}

type ready struct{}

func volumeMeasuringService(
   meter metronome.VolumeMeter,
   terminationRequestChannel <-chan terminate) (
   volumeChannel <-chan int,
   firstVolumeReadyChannel <-chan ready) {

   channel := make(chan int, 5)
   readyChannel := make(chan ready)

   go func() {
      volume := meter()
      close(readyChannel)
      channel <- volume

      for !isTerminating(terminationRequestChannel) {
         volume := meter()
         channel <- volume
      }
      close(channel)
   }()

   return channel, readyChannel
}

type complete struct{}

func getMostRecentMessage(intChannel <-chan int, defaultInt int) (mostRecentOrDefaultIfNone int) {
   mostRecentOrDefaultIfNone = defaultInt
   for {
      select {
      case newVolume, ok := <-intChannel:
         if ok {
            mostRecentOrDefaultIfNone = newVolume
         }
      default:
         return
      }
   }
}

func beatPerformingService(
   performer metronome.BeatPerformer,
   beatCountChannel <-chan int,
   timingChannel <-chan timing,
   volumeChannel <-chan int,
   terminationRequestChannel chan<- terminate,
) (beatPerformanceCompleteChannel <-chan complete) {

   channel := make(chan complete)

   go func() {
      var volume int
      for beatCount := range beatCountChannel {
         <-timingChannel

         volume = getMostRecentMessage(volumeChannel, volume)

         performer(beatCount, volume)
      }

      close(terminationRequestChannel)

      for range timingChannel {
      }

      for range volumeChannel {
      }

      close(channel)
   }()

   return channel
}

func asynchronousWithChannelsServices(
   volumeMeter metronome.VolumeMeter,
   numberOfBeats int,
   _ track.Tracker) metronome.Metronome {

   return func(bpm param.Bpm, performer metronome.BeatPerformer) {
      terminationRequestChannel := make(chan terminate)

      beatCountChannel := countingService(numberOfBeats)
      volumeChannel, firstVolumeReadyChannel := volumeMeasuringService(volumeMeter, terminationRequestChannel)
      timingChannel := timingService(bpm.Interval(), firstVolumeReadyChannel, terminationRequestChannel)

      beatPerformanceCompleteChannel := beatPerformingService(
         performer,
         beatCountChannel,
         timingChannel,
         volumeChannel,
         terminationRequestChannel,
      )

      for range beatPerformanceCompleteChannel {
      }

   }
}

func AsynchronousWithChannelsServices(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
   demonstrateScenario("async_with_channels_services", generator.RandomisedWithProportionalCapsDuration, bpm, numberOfBeats, false, metas, eventsChannel,
      func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
         asynchronousWithChannelsServices(volumeMeter, numberOfBeats, tracker)(bpm, performer)
         return chart.Meta{
            Description: "Asynchronous metronome with <b>channels</b> communicating state changes. All asynchronous. " +
               "Execution time of volume measurement is <b>random</b>. ",
            Groups: defaultGroups,
         }
      },
   )
}
