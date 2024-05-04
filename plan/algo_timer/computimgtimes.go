package algo_timer

import (
	"fmt"
	"time"
)

// NewTimer
func NewTimer() *Timer {
	timer := &Timer{mission: true}
	return timer
}

// Close the timer
func (t *Timer) TimerClose() {
	t.mission = false
}

// Open the timer
func (t *Timer) TimerOpen() {
	t.mission = true
}

// Start the timer
func (t *Timer) TimerStart() {
	if t.mission {
		t.start = time.Duration(time.Now().UnixMicro())
	}
}

// Stop the timer
func (t *Timer) TimerStop() {
	if t.mission {
		t.end = time.Duration(time.Now().UnixMicro())
		t.total += t.end - t.start
	}
}

// Stop the timer
func (t *Timer) TimerEnd() {
	if t.mission {
		t.end = time.Duration(time.Now().UnixMicro())
		t.total += t.end - t.start

		// Mix limit 1000ms
		if t.total >= time.Duration(1000)*time.Millisecond {
			t.total = time.Duration(1000) * time.Millisecond
		}
	}
	t.mission = false
}

func (t *Timer) TimerMax() {
	if t.mission {
		t.total = time.Duration(1000) * time.Millisecond
		t.mission = false
	}
}

// Export data
func (t *Timer) TimerExportData() {
	fmt.Printf(" %v \n", t.total)
}

// Output data
func (t *Timer) TimerOutputData() time.Duration {
	return t.total
}

// Merge the timer
func (t1 *Timer) TimerMerge(t2 *Timer) {
	if t1.mission {
		t1.total = t1.total + t2.total
		t1.mission = false

		// Mix limit 1000ms
		if t1.total >= time.Duration(1000)*time.Millisecond {
			t1.total = time.Duration(1000) * time.Millisecond
		}
	}

}
