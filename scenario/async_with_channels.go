package scenario

import (
   "github.com/outo/metronome/generator"
   "github.com/outo/metronome/track"
   "github.com/outo/metronome/param"
   "github.com/outo/metronome/chart"
   "github.com/outo/metronome"
   "time"
)

func asynchronousWithChannels(volumeMeter metronome.VolumeMeter, numberOfBeats int, _ track.Tracker) metronome.Metronome {
   return func(bpm param.Bpm, performer metronome.BeatPerformer) {

      // channel will carry requests for volume from main loop to the volume measuring Goroutine
      // I had to add buffer of one as, straight from start, I need to send two requests without
      // receiver attached (which would block)
      volumeRequestsChannel := make(chan bool, 1)

      // channel will deliver volume measurements from volume measuring Goroutine to the main loop
      // no need for buffering as main loop will swiftly receive all messages
      volumeResponsesChannel := make(chan int)

      // Goroutine which drains volume requests channel, measures
      // volume on each request, then sends the volume
      // in response message
      go func() {
         for range volumeRequestsChannel {
            // request arrived

            // blocking volume measurement (but does not impact main timing
            // loop as this operation runs in Goroutine)
            volume := volumeMeter()

            // send the latest measurement as response
            volumeResponsesChannel <- volume
         }
         // at this point we know that volume requests channel
         // is closed and empty (because for-range loop ended)

         // close volume response channel to indicate this Goroutine is ending
         close(volumeResponsesChannel)
      }()

      // request first volume measurement, so that we know when to start the timer
      volumeRequestsChannel <- true

      // request next volume measurement
      volumeRequestsChannel <- true

      // once we get the first volume measurement the timer can start
      // without this the time span between the first two beats (only)
      // would likely be longer than desired interval
      volume := <-volumeResponsesChannel

      // starting timer right after first volume arrived ensures that the next timing event
      // will have at least this volume to work with (otherwise it would have to wait for some volume or default it)
      ticker := time.NewTicker(bpm.Interval())
      defer ticker.Stop()

      // first beat performed synchronously to avoid complicating the loop's code
      performer(0, volume)

      // for-loop until we have performed numberOfBeats - 1 (one already done above),
      // notice this loop does not increment beatCount as the beats don't happen on each iteration
      // this loop resembles more of a "while loop"
      for beatCount := 1; beatCount < numberOfBeats; {
         // non-blocking select which will try to retrieve timing message or volume,
         // in case none available, it will iterate again immediately, until one of the
         // messages is available
         select {
         case <-ticker.C:
            // we got timing message back from Ticker, it's time to perform a beat
            performer(beatCount, volume)
            // and increment the count
            beatCount++
         case volume = <-volumeResponsesChannel:
            // fresh volume value arrived, overwrite what's in the volume variable
            // and request next volume
            volumeRequestsChannel <- true
         default:
         }
      }
      //indicate to volume measuring Goroutine that there will be no more requests
      // that will end for-reach loop in Goroutine
      close(volumeRequestsChannel)

      //wait for last volume response to be delivered and volume response channel closed
      for range volumeResponsesChannel {
      }

      // Goroutine finished, both channels are drained and empty
      // we can finish the program gracefully now
   }
}

func AsynchronousWithChannels(bpm param.Bpm, numberOfBeats int, metas map[track.ExecutionId]chart.Meta, eventsChannel chan track.Event) {
   demonstrateScenario("async_with_channels", generator.RandomisedWithProportionalCapsDuration, bpm, numberOfBeats, false, metas, eventsChannel,
      func(volumeMeter metronome.VolumeMeter, tracker track.Tracker, performer metronome.BeatPerformer) chart.Meta {
         asynchronousWithChannels(volumeMeter, numberOfBeats, tracker)(bpm, performer)
         return chart.Meta{
            Description: "Asynchronous metronome with <b>channels</b> communicating state changes. " +
               "Execution time of volume measurement is <b>random</b>. ",
            Groups: defaultGroups,
         }
      },
   )
}
