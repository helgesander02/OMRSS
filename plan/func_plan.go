package plan

import (
	"src/internal/random"
	"src/plan/algo"
	"src/plan/routes"
)

func FillRNG(r *random.Generator) {
	routes.SetRNG(r)
	algo.SetRNG(r)
}
