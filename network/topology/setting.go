package topology

import "src/pkg/random"

var rng *random.Generator

func FillRNG(r *random.Generator) {
	rng = r
}
