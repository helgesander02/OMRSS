package random

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Generator struct {
	rng *rand.Rand
}

func New(seed int64) *Generator {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	source := rand.NewPCG(uint64(seed), uint64(seed))

	return &Generator{
		rng: rand.New(source),
	}
}

func (g *Generator) IntN(n int) int {
	return g.rng.IntN(n)
}

func (g *Generator) Int64N(n int64) int64 {
	return g.rng.Int64N(n)
}

func (g *Generator) Float64() float64 {
	return g.rng.Float64()
}

func (g *Generator) Shuffle(n int, swap func(i, j int)) {
	g.rng.Shuffle(n, swap)
}

func SelectRandom[T any](g *Generator, slice []T) (T, error) {
	var zero T

	if len(slice) == 0 {
		return zero, fmt.Errorf("cannot select from empty slice")
	}

	return slice[g.IntN(len(slice))], nil
}

func (g *Generator) SelectRandomIndex(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("n must be positive, got: %d", n)
	}

	return g.IntN(n), nil
}

func SelectRandomN[T any](g *Generator, slice []T, n int) ([]T, error) {

	if n > len(slice) {
		return nil, fmt.Errorf(
			"cannot select %d elements from slice of length %d",
			n,
			len(slice),
		)
	}

	if n <= 0 {
		return []T{}, nil
	}

	indices := make([]int, len(slice))

	for i := range indices {
		indices[i] = i
	}

	g.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	result := make([]T, n)

	for i := 0; i < n; i++ {
		result[i] = slice[indices[i]]
	}

	return result, nil
}

func (g *Generator) SelectRandomNIndices(max, n int) ([]int, error) {

	if n > max {
		return nil, fmt.Errorf(
			"cannot select %d indices from range [0, %d)",
			n,
			max,
		)
	}

	if n <= 0 {
		return []int{}, nil
	}

	indices := make([]int, max)

	for i := range indices {
		indices[i] = i
	}

	g.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	return indices[:n], nil
}

func (g *Generator) WeightedSelect(weights []float64) (int, error) {

	if len(weights) == 0 {
		return 0, fmt.Errorf("weights slice is empty")
	}

	sum := 0.0

	for _, w := range weights {
		sum += w
	}

	if sum <= 0 {
		return 0, fmt.Errorf(
			"sum of weights must be positive, got: %f",
			sum,
		)
	}

	r := g.Float64() * sum
	cumulative := 0.0

	for i, w := range weights {
		cumulative += w

		if r < cumulative {
			return i, nil
		}
	}

	// Fallback due to floating-point precision
	return len(weights) - 1, nil
}
