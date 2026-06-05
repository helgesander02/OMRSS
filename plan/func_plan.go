package plan

import (
	"src/pkg/random"
	"src/plan/algo"
	"src/plan/routes"
)

func FillRNG(r *random.Generator) {
	routes.SetRNG(r)
	algo.SetRNG(r)
}
