package timingwheel

import (
	"fmt"
	"testing"
	"time"
)

func TestTimingWheel_AfterFunc(t *testing.T) {
	tw := NewTimingWheel(time.Millisecond, 10)
	tw.Start()
	defer tw.Stop()

	durations := []time.Duration{
		2 * time.Millisecond,
		2 * time.Millisecond,
		2 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		3 * time.Millisecond,
		10 * time.Millisecond,
		15 * time.Millisecond,
	}

	exitC := make(chan time.Time)

	for _, d := range durations {
		tw.AfterFunc(d, func() {
			exitC <- time.Now().UTC()
		})
	}

	for t := range exitC {
		fmt.Println(t)
	}
}
