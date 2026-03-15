package topology

import "src/internal/random"

const (
	CANNodeCount = 5
)

var rng *random.Generator

func FillRNG(r *random.Generator) {
	rng = r
}
