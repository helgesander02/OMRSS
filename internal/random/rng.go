package random

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// Generator wraps rand.Rand for dependency injection and testing
type Generator struct {
	rng *rand.Rand
}

// New creates a new Generator with the given seed
// If seed is 0, uses current time as seed
func New(seed int64) *Generator {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	// Use PCG algorithm for better randomness and reproducibility
	source := rand.NewPCG(uint64(seed), uint64(seed))
	return &Generator{
		rng: rand.New(source),
	}
}

// IntN returns a random integer in [0, n)
// Panics if n <= 0
func (g *Generator) IntN(n int) int {
	return g.rng.IntN(n)
}

// Int64N returns a random int64 in [0, n)
// Panics if n <= 0
func (g *Generator) Int64N(n int64) int64 {
	return g.rng.Int64N(n)
}

// Float64 returns a random float64 in [0.0, 1.0)
func (g *Generator) Float64() float64 {
	return g.rng.Float64()
}

// Shuffle randomizes the order of elements
func (g *Generator) Shuffle(n int, swap func(i, j int)) {
	g.rng.Shuffle(n, swap)
}

// SelectRandom selects a random element from a slice
func SelectRandom[T any](g *Generator, slice []T) (T, error) {
	var zero T
	if len(slice) == 0 {
		return zero, fmt.Errorf("cannot select from empty slice")
	}
	return slice[g.IntN(len(slice))], nil
}

// SelectRandomIndex returns a random index from [0, n)
func (g *Generator) SelectRandomIndex(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("n must be positive, got: %d", n)
	}
	return g.IntN(n), nil
}

// SelectRandomN selects n random unique elements from a slice
func SelectRandomN[T any](g *Generator, slice []T, n int) ([]T, error) {
	if n > len(slice) {
		return nil, fmt.Errorf("cannot select %d elements from slice of length %d", n, len(slice))
	}

	if n <= 0 {
		return []T{}, nil
	}

	// Create a copy of indices
	indices := make([]int, len(slice))
	for i := range indices {
		indices[i] = i
	}

	// Shuffle and select first n
	g.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = slice[indices[i]]
	}

	return result, nil
}

// SelectRandomNIndices selects n random unique indices from [0, max)
func (g *Generator) SelectRandomNIndices(max, n int) ([]int, error) {
	if n > max {
		return nil, fmt.Errorf("cannot select %d indices from range [0, %d)", n, max)
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

// WeightedSelect selects an index based on weighted probability
// weights should sum to approximately 1.0, but will be normalized internally
func (g *Generator) WeightedSelect(weights []float64) (int, error) {
	if len(weights) == 0 {
		return 0, fmt.Errorf("weights slice is empty")
	}

	// Normalize weights
	sum := 0.0
	for _, w := range weights {
		sum += w
	}

	if sum <= 0 {
		return 0, fmt.Errorf("sum of weights must be positive, got: %f", sum)
	}

	r := g.Float64() * sum
	cumulative := 0.0

	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return i, nil
		}
	}

	// Fallback to last index due to floating point precision
	return len(weights) - 1, nil
}
