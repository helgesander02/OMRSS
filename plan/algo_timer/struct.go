package algo_timer

import (
	"time"
)

// Timer structure
type Timer struct {
	mission bool
	total   time.Duration
	start   time.Duration
	end     time.Duration
}
