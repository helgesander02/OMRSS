package schedule

import (
	"fmt"
	"time"
)

// Timer structure
type Timer struct {
	mission bool
	total   time.Duration
	start   time.Duration
	end     time.Duration
}

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
	t.start = time.Duration(time.Now().UnixMicro())
}

// Stop the timer
func (t *Timer) TimerStop() {
	t.end = time.Duration(time.Now().UnixMicro())
	t.total += t.end - t.start
}

// Stop the timer
func (t *Timer) TimerEnd() {
	t.end = time.Duration(time.Now().UnixMicro())
	t.total += t.end - t.start
}

// Mix limit 1000ms
func (t *Timer) Exceededlimit() {
	t.total = time.Duration(1000) * time.Millisecond
}

// Export data
func (t *Timer) TimerExportData() {
	fmt.Printf(" %v \n", t.total)
}

// Merge the timer
func (t1 *Timer) TimerMerge(t2 *Timer) {
	fmt.Printf(" %v \n", t1.total+t2.total)
}
