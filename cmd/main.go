package main

import (
   "github.com/outo/metronome/param"
   . "github.com/outo/metronome/track"
   "github.com/outo/metronome/chart"
   "github.com/outo/metronome/scenario"
   "runtime"
   "fmt"
)

func main() {
   //fast rate so that the demo executes swiftly
   //setting it too high (depends on machine) will
   // skew the results due to under performing CPU
   bpm := param.Bpm(6000)

   //enough beats to illustrate the problems and solutions
   const numberOfBeats = 15

   chartMeta := make(map[ExecutionId]chart.Meta)

   //buffered channel to accommodate all events
   //because I'm reading the events only once all the demos complete,
   //not providing sufficient head-room would cause deadlock
   const numberOfEventsPerChart = 3 // increase if you alter the chart to contain more events (channel deadlock otherwise!)
   const numberOfCharts = 14        /*some cases produce more than 1 */
   eventsChannel := make(chan Event, numberOfBeats*numberOfEventsPerChart*numberOfCharts)

   //emphasise the impact of badly written concurrent code by enabling true parallelism on machines with multiple cores
   runtime.GOMAXPROCS(runtime.NumCPU())

   scenario.SynchronousSleepWithPreMeter(bpm, numberOfBeats, chartMeta, eventsChannel)
   scenario.SynchronousConstantDelay(bpm, numberOfBeats, chartMeta, eventsChannel)
   scenario.SynchronousAdaptiveDelay(bpm, numberOfBeats, chartMeta, eventsChannel)
   scenario.SynchronousTicker(bpm, numberOfBeats, chartMeta, eventsChannel)
   scenario.SynchronousTickerFireImmediately(bpm, numberOfBeats, chartMeta, eventsChannel)
   scenario.AsynchronousButNaive(bpm, numberOfBeats, chartMeta, eventsChannel)
   scenario.AsynchronousStillNaiveWithWaitGroup(bpm, numberOfBeats, chartMeta, eventsChannel)
   scenario.AsynchronousWithAtomic(bpm, numberOfBeats, chartMeta, eventsChannel)
   scenario.AsynchronousWithMutex(bpm, numberOfBeats, chartMeta, eventsChannel)
   scenario.AsynchronousWithChannels(bpm, numberOfBeats, chartMeta, eventsChannel)
   scenario.AsynchronousWithChannelsServices(bpm, numberOfBeats, chartMeta, eventsChannel)

   close(eventsChannel)

   for i:=0; i < 30; i++ {
      fmt.Println()
   }

   err := chart.AllEvents(chartMeta, eventsChannel)
   if err != nil {
      panic(err)
   }
}
