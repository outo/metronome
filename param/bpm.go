package param

import (
	"time"
	"errors"
)

type Bpm int

func (bpm Bpm) Interval() time.Duration {
	return time.Duration(60 / float32(bpm) * float32(time.Second))
}

func New(value int) (bpm Bpm, err error) {
	if value <= 0 || value > 6000 {
		err = errors.New("bpm has to within (0, 6000]")
		return
	}

	bpm = Bpm(value)
	return
}