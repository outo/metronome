package generator

import (
	"time"
	"math/rand"
)

type Duration func(time.Duration) time.Duration
type PreconfiguredDuration func() (time.Duration)

var (
	//it is not relevant what this time exactly is, as long as we have some recent, static time stored here
	someTime = time.Now()

	ShortDuration Duration = func(duration time.Duration) time.Duration {
		return time.Duration(float64(duration.Nanoseconds()) * 0.5)
	}

	RandomisedWithProportionalCapsDuration Duration = func(duration time.Duration) time.Duration {
		minFactor := float32(0.2)
		maxFactor := float32(2.5)
		return time.Duration(float32(duration)*minFactor + (float32(duration) * (maxFactor - minFactor) * rand.Float32()))
	}

	TriangularPatternDuration Duration = func(duration time.Duration) time.Duration {

		minY := time.Duration(float64(duration.Nanoseconds()) * 0.2)
		maxY := time.Duration(float64(duration.Nanoseconds()) * 2.5)
		period := duration * 20
		timeDiff := time.Now().Sub(someTime)
		x := float64(timeDiff.Nanoseconds() % period.Nanoseconds() / duration.Nanoseconds())

		var value int64
		if x < 10 {
			value = int64(x)
		} else {
			value = int64(19 - x)
		}

		y := minY + time.Duration(value*(maxY - minY).Nanoseconds()/10)

		return y
	}
)
