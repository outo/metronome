package track

import (
	"time"
	"fmt"
)

type ExecutionId string

type Event struct {
	ExecutionId
	Category       string
	Content        string
	Description    string
	TimeSinceStart time.Duration
	Duration       time.Duration
	ChartItemClassName,
	ChartItemType string
}

func (e Event) String() (s string) {
	appendIfTrue := func(condition bool, text string) {
		if condition {
			s += text
		}
	}
	appendIfTrue(e.Category != "", fmt.Sprintf("Category: %s ", e.Category))
	appendIfTrue(e.Content != "", fmt.Sprintf("Content: %s ", e.Content))
	appendIfTrue(e.Description != "", fmt.Sprintf("Description: %s ", e.Description))
	appendIfTrue(true, fmt.Sprintf("TimeSinceStart: %v ", e.TimeSinceStart))
	appendIfTrue(e.Duration != 0, fmt.Sprintf("Duration: %v ", e.Duration))
	appendIfTrue(e.ChartItemClassName != "", fmt.Sprintf("ChartItemClassName: %s ", e.ChartItemClassName))
	appendIfTrue(e.ChartItemType != "", fmt.Sprintf("ChartItemType: %s ", e.ChartItemType))
	return "track: " + s
}

type Tracker func(event Event)

func New(executionId ExecutionId, channel chan Event) Tracker {
	thisTrackerStartedAt := time.Now()
	return func(event Event) {
		event.ExecutionId = executionId
		if event.TimeSinceStart == 0 {
			if event.Duration != 0 {
				event.TimeSinceStart = time.Since(thisTrackerStartedAt) - event.Duration
			} else {
				event.TimeSinceStart = time.Since(thisTrackerStartedAt)
			}
		}

		// this is to recover from panic due to demo of "naive asynchronous" implementation (if executed as last demo)
		// normally the producer on a channel would know when not to send more messages
		// but that demo is demonstrating an ignorant approach which does not care
		defer func() {
			if rc := recover(); rc != nil {
				fmt.Println("Recovered from panic:", rc)
			}
		}()

		channel <- event
	}
}

func VerticalMarker(tracker Tracker) {
	tracker(Event{Category: "custom_time"})
}
