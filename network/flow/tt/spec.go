package tt

import (
	"src/internal/random"
)

// Global RNG instance (will be set by network package)
var rng *random.Generator

// SetRNG sets the random number generator for this package
func SetRNG(r *random.Generator) {
	rng = r
}

func config_TSN_Stream() *TSN {
	tPeriod, tDatasize := random_TSN()
	tsn := new_TSN(tPeriod, tDatasize)

	return tsn
}

func config_AVB_Stream() *AVB {
	aDatasize := random_AVB()
	avb := new_AVB(aDatasize)

	return avb
}

func random_TSN() (int, float64) {
	tsnPeriodArr := []int{100, 500, 1000, 1500, 2000}
	tsnDatasizeArr := []float64{30., 40., 50., 60., 70., 80., 90., 100.}

	periodIdx := rng.IntN(len(tsnPeriodArr))
	datasizeIdx := rng.IntN(len(tsnDatasizeArr))

	return tsnPeriodArr[periodIdx], tsnDatasizeArr[datasizeIdx]
}

func random_AVB() float64 {
	avbDatasizeArr := []float64{1000., 1100., 1200., 1300., 1400., 1500.}
	datasizeIdx := rng.IntN(len(avbDatasizeArr))

	return avbDatasizeArr[datasizeIdx]
}

func random_TT_Devices_For_Tree(Nnode int) (int, []int) {
	// Talker
	sourceIdx := rng.IntN(Nnode)

	// Listener - all nodes except source
	destinations := []int{}
	for i := 0; i < Nnode; i++ {
		if i != sourceIdx {
			destinations = append(destinations, i+2000)
		}
	}

	// Calculate number of destinations: random(0, 1) + (Nnode-1-3) + 3 = random(Nnode-1, Nnode)
	// This ensures at least Nnode-1 destinations and at most Nnode destinations
	baseOffset := rng.IntN(2) // 0 or 1
	maxRange := Nnode - 1 - 3 // Nnode - 4
	if maxRange < 0 {
		maxRange = 0
	}
	randomOffset := 0
	if maxRange > 0 {
		randomOffset = rng.IntN(maxRange + 1) // 0 to maxRange inclusive
	}
	numDestinations := baseOffset + randomOffset + 3

	// Ensure numDestinations doesn't exceed available destinations
	if numDestinations > len(destinations) {
		numDestinations = len(destinations)
	}

	// Randomly select destinations without replacement
	selectedDestinations := []int{}
	destCopy := make([]int, len(destinations))
	copy(destCopy, destinations)

	for i := 0; i < numDestinations; i++ {
		randIndex := rng.IntN(len(destCopy))
		selectedDestinations = append(selectedDestinations, destCopy[randIndex])
		// Remove selected element
		destCopy = append(destCopy[:randIndex], destCopy[randIndex+1:]...)
	}

	return sourceIdx + 1000, selectedDestinations // source id: 1000+id  destination id: 2000+id
}
