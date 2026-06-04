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

// Global parameter storage for TSN flows
var (
	tsnPeriods   []int
	tsnDataSizes []float64
)

// Global parameter storage for AVB flows
var (
	avbPeriod    int
	avbDeadline  int
	avbDataSizes []float64
)

// SetTSNParams sets the TSN flow parameters from config
func SetTSNParams(periods []int, dataSizes []float64) {
	tsnPeriods = periods
	tsnDataSizes = dataSizes
}

// SetAVBParams sets the AVB flow parameters from config
func SetAVBParams(period, deadline int, dataSizes []float64) {
	avbPeriod = period
	avbDeadline = deadline
	avbDataSizes = dataSizes
}

func configTSNFrame() *TSN {
	tPeriod, tDatasize := randomTSN()
	tsn := newTSN(tPeriod, tDatasize)

	return tsn
}

func configAVBFrame() *AVB {
	aDatasize := randomAVB()
	avb := newAVB(aDatasize)

	return avb
}

func randomTSN() (int, float64) {
	// Use config values if set, otherwise use defaults
	periods := tsnPeriods
	if len(periods) == 0 {
		periods = []int{100, 500, 1000, 1500, 2000}
	}

	dataSizes := tsnDataSizes
	if len(dataSizes) == 0 {
		dataSizes = []float64{30., 40., 50., 60., 70., 80., 90., 100.}
	}

	periodIdx := rng.IntN(len(periods))
	datasizeIdx := rng.IntN(len(dataSizes))

	return periods[periodIdx], dataSizes[datasizeIdx]
}

func randomAVB() float64 {
	// Use config values if set, otherwise use defaults
	dataSizes := avbDataSizes
	if len(dataSizes) == 0 {
		dataSizes = []float64{1000., 1100., 1200., 1300., 1400., 1500.}
	}

	datasizeIdx := rng.IntN(len(dataSizes))

	return dataSizes[datasizeIdx]
}

func randomTTDevicesForTree(nnode int) (int, []int) {
	// Talker
	sourceIdx := rng.IntN(nnode)

	// Listener - all nodes except source
	destinations := []int{}
	for i := 0; i < nnode; i++ {
		if i != sourceIdx {
			destinations = append(destinations, i+2000)
		}
	}

	// Calculate number of destinations: random(0, 1) + (nnode-1-3) + 3 = random(nnode-1, nnode)
	// This ensures at least nnode-1 destinations and at most nnode destinations
	baseOffset := rng.IntN(2) // 0 or 1
	maxRange := nnode - 1 - 3 // nnode - 4
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
