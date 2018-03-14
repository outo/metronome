package scenario

import (
	"github.com/outo/metronome/generator"
	"github.com/outo/metronome/track"
	"github.com/outo/metronome/chart"
	"github.com/outo/metronome/param"
	"fmt"
	"github.com/outo/metronome/volumemeter"
	"strings"
	"github.com/outo/metronome"
)

type MetronomeFactory func(numberOfBeats int, delaySimulator generator.Duration, tracker track.Tracker) metronome.Metronome

var (
	ghostGroup  = chart.Group{Id: "0", Content: "Simulated frequency <br>of ideal beats", DefaultItemClassName: "ghost"}
	volumeGroup = chart.Group{Id: "1", Content: "Volume measurement", DefaultItemClassName: "volume", DefaultItemType: "background"}
	beatGroup   = chart.Group{Id: "2", Content: "Actual timing of beats", DefaultItemClassName: "beat"}

	defaultGroups = map[string]chart.Group{
		"ghost":  ghostGroup,
		"volume": volumeGroup,
		"beat":   beatGroup,
	}
)

func demonstrateScenario(
	id track.ExecutionId,
	delaySimulator generator.Duration,
	bpm param.Bpm,
	numberOfBeats int,
	printBeatCount bool,
	metas map[track.ExecutionId]chart.Meta,
	eventsChannel chan track.Event,
	runScenario func(metronome.VolumeMeter, track.Tracker, metronome.BeatPerformer) chart.Meta) {

	fmt.Println("Recording execution of", id)

	executionId := track.ExecutionId(id)
	tracker := track.New(executionId, eventsChannel)
	chart.ApplyGhosting(bpm, numberOfBeats, tracker)
	var performer metronome.BeatPerformer

	if printBeatCount {
		performer = metronome.CreateBeatPerformerWithTrackingWithCount(tracker)
	} else {
		performer = metronome.CreateBeatPerformerWithTrackingWithIOContent(tracker)
	}
	volumeMeter := volumemeter.New(bpm.Interval(), delaySimulator, tracker)

	meta := runScenario(volumeMeter, tracker, performer)
	if strings.Contains(string(id), "async") {
		track.VerticalMarker(tracker)
	}

	length := len(metas)
	meta.Order = length
	meta.Description = fmt.Sprintf("%d. %s", length+1, meta.Description)
	metas[executionId] = meta
}
