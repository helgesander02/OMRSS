package tt

import "src/pkg/random"

// random variables
var rng *random.Generator

func SetRNG(r *random.Generator) {
	rng = r
}

// parameters from config
var (
	tsnPeriods   []int
	tsnDataSizes []float64

	avbPeriod    int
	avbDeadline  int
	avbDataSizes []float64
)

func SetTSNParams(periods []int, dataSizes []float64) {
	tsnPeriods = periods
	tsnDataSizes = dataSizes
}

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
	periodIdx := rng.IntN(len(tsnPeriods))
	datasizeIdx := rng.IntN(len(tsnDataSizes))

	return tsnPeriods[periodIdx], tsnDataSizes[datasizeIdx]
}

func randomAVB() float64 {
	datasizeIdx := rng.IntN(len(avbDataSizes))

	return avbDataSizes[datasizeIdx]
}

func randomTTDevices(nnode int, mode string) (int, []int) {
	// Talker
	sourceIdx := rng.IntN(nnode)
	selectedSource := sourceIdx + 1000

	// Listener - all nodes except source
	destinations := []int{}
	for i := range nnode {
		if i != sourceIdx {
			destinations = append(destinations, i+2000)
		}
	}

	selectedDestinations := []int{}

	destCopy := make([]int, len(destinations))
	copy(destCopy, destinations)

	// Path: select only one destination
	if mode == "path" {
		randIndex := rng.IntN(len(destCopy))
		selectedDestinations = append(selectedDestinations, destCopy[randIndex])
		return selectedSource, selectedDestinations
	}

	// Tree: select multiple destinations
	baseOffset := rng.IntN(2) // 0 or 1
	maxRange := nnode - 4
	if maxRange < 0 {
		maxRange = 0
	}

	randomOffset := 0
	if maxRange > 0 {
		randomOffset = rng.IntN(maxRange + 1)
	}

	numDestinations := baseOffset + randomOffset + 3

	if numDestinations > len(destCopy) {
		numDestinations = len(destCopy)
	}

	for range numDestinations {
		randIndex := rng.IntN(len(destCopy))
		selectedDestinations = append(selectedDestinations, destCopy[randIndex])
		destCopy = append(destCopy[:randIndex], destCopy[randIndex+1:]...)
	}

	return selectedSource, selectedDestinations // source id: 1000+id  destination id: 2000+id
}
