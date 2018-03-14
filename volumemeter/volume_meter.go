package volumemeter

import (
	"time"
	"github.com/outo/metronome/track"
	"github.com/outo/metronome/generator"
	"math/rand"
	"github.com/outo/metronome"
	"strconv"
)


func New(
	baselineDuration time.Duration,
	delaySimulator generator.Duration,
	tracker track.Tracker) metronome.VolumeMeter {

	return func() (volume int) {

		//use delay simulator to generate a delay corresponding to some baseline
		// baseline is needed so that the delay is not 1s for an application working at 1ns resolution
		simulatedDelay := delaySimulator(baselineDuration)

		//pretend the execution of "volume measuring" takes more time than it should
		time.Sleep(simulatedDelay)

		// Obviously the volume is random (well, pseudo-random) just for the purpose of the demo.
		// In real-life, if the volume was to change so instantly between extremities there
		// would be little point trying to measure it to produce beats.
		volume = rand.Intn(101)

		event := track.Event{
			Category: "volume",
			Content:  strconv.Itoa(volume),
			Duration: simulatedDelay,
		}

		tracker(event)

		return
	}
}
