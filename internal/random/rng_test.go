package random

import (
	"testing"
)

func TestRNG_Reproducibility(t *testing.T) {
	seed := int64(12345)

	// First run
	rng1 := New(seed)
	results1 := make([]int, 100)
	for i := 0; i < 100; i++ {
		results1[i] = rng1.IntN(1000)
	}

	// Second run with same seed
	rng2 := New(seed)
	results2 := make([]int, 100)
	for i := 0; i < 100; i++ {
		results2[i] = rng2.IntN(1000)
	}

	// Results should be identical
	for i := 0; i < 100; i++ {
		if results1[i] != results2[i] {
			t.Errorf("Results differ at index %d: %d != %d", i, results1[i], results2[i])
		}
	}
}

func TestRNG_DifferentSeeds(t *testing.T) {
	rng1 := New(12345)
	rng2 := New(54321)

	// Generate some numbers
	same := 0
	total := 100
	for i := 0; i < total; i++ {
		v1 := rng1.IntN(1000)
		v2 := rng2.IntN(1000)
		if v1 == v2 {
			same++
		}
	}

	// Should have very few collisions (statistically)
	if same > total/2 {
		t.Errorf("Too many collisions with different seeds: %d/%d", same, total)
	}
}

func TestSelectRandom(t *testing.T) {
	rng := New(42)
	slice := []string{"a", "b", "c", "d", "e"}

	// Test normal case
	for i := 0; i < 10; i++ {
		elem, err := SelectRandom(rng, slice)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify element is in slice
		found := false
		for _, v := range slice {
			if v == elem {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Selected element %s not in original slice", elem)
		}
	}

	// Test empty slice
	emptySlice := []string{}
	_, err := SelectRandom(rng, emptySlice)
	if err == nil {
		t.Error("Expected error for empty slice, got nil")
	}
}

func TestSelectRandomN(t *testing.T) {
	rng := New(42)
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Test normal case
	selected, err := SelectRandomN(rng, slice, 5)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(selected) != 5 {
		t.Errorf("Expected 5 elements, got %d", len(selected))
	}

	// Verify uniqueness
	seen := make(map[int]bool)
	for _, v := range selected {
		if seen[v] {
			t.Errorf("Duplicate element: %d", v)
		}
		seen[v] = true
	}

	// Test edge case: select all elements
	all, err := SelectRandomN(rng, slice, len(slice))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(all) != len(slice) {
		t.Errorf("Expected %d elements, got %d", len(slice), len(all))
	}

	// Test error case: select more than available
	_, err = SelectRandomN(rng, slice, len(slice)+1)
	if err == nil {
		t.Error("Expected error when selecting more elements than available")
	}
}

func TestWeightedSelect(t *testing.T) {
	rng := New(42)

	// Test with equal weights - should be approximately uniform
	weights := []float64{1.0, 1.0, 1.0, 1.0}
	counts := make([]int, len(weights))

	iterations := 10000
	for i := 0; i < iterations; i++ {
		idx, err := rng.WeightedSelect(weights)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		counts[idx]++
	}

	// Each should be approximately 25% (2500 ± some variance)
	expected := iterations / len(weights)
	tolerance := expected / 5 // 20% tolerance
	for i, count := range counts {
		if count < expected-tolerance || count > expected+tolerance {
			t.Logf("Warning: Count for index %d is %d, expected around %d", i, count, expected)
		}
	}

	// Test with heavily skewed weights
	skewedWeights := []float64{10.0, 1.0, 1.0, 1.0}
	skewedCounts := make([]int, len(skewedWeights))

	for i := 0; i < iterations; i++ {
		idx, err := rng.WeightedSelect(skewedWeights)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		skewedCounts[idx]++
	}

	// First element should be selected much more often (~77% of the time)
	if skewedCounts[0] < iterations/2 {
		t.Errorf("First element (weight 10.0) should be selected more often, got %d/%d", skewedCounts[0], iterations)
	}

	// Test error cases
	_, err := rng.WeightedSelect([]float64{})
	if err == nil {
		t.Error("Expected error for empty weights")
	}

	_, err = rng.WeightedSelect([]float64{0, 0, 0})
	if err == nil {
		t.Error("Expected error for all-zero weights")
	}
}

func BenchmarkIntN(b *testing.B) {
	rng := New(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rng.IntN(1000)
	}
}

func BenchmarkSelectRandom(b *testing.B) {
	rng := New(42)
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SelectRandom(rng, slice)
	}
}

func BenchmarkWeightedSelect(b *testing.B) {
	rng := New(42)
	weights := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = rng.WeightedSelect(weights)
	}
}
