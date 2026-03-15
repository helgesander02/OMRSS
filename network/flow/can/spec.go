package can

import (
	"src/internal/random"
)

// Global RNG instance (will be set by network package)
var rng *random.Generator

// Global CAN parameters (will be set from config)
var importantCANParams struct {
	period   int
	deadline int
	dataSize float64
}

var unimportantCANParams struct {
	periods   []int
	deadlines []int
	dataSize  float64
}

// SetRNG sets the random number generator for this package
func SetRNG(r *random.Generator) {
	rng = r
}

// SetImportantCANParams sets the parameters for important CAN flows from config
func SetImportantCANParams(period int, deadline int, dataSize float64) {
	importantCANParams.period = period
	importantCANParams.deadline = deadline
	importantCANParams.dataSize = dataSize
}

// SetUnimportantCANParams sets the parameters for unimportant CAN flows from config
func SetUnimportantCANParams(periods []int, deadlines []int, dataSize float64) {
	unimportantCANParams.periods = periods
	unimportantCANParams.deadlines = deadlines
	unimportantCANParams.dataSize = dataSize
}

func configImportantCANStream() *importantCAN {
	// Use configured parameters if available, otherwise use defaults
	if importantCANParams.period > 0 {
		return newImportantCANWithParams(
			importantCANParams.period,
			importantCANParams.deadline,
			importantCANParams.dataSize,
		)
	}
	return newImportantCAN()
}

func configUnimportantCANStream() *unimportantCAN {
	ucPeriod, ucDeadline := randomUnimportantCAN()

	// Use configured data size if available
	dataSize := 16.0
	if unimportantCANParams.dataSize > 0 {
		dataSize = unimportantCANParams.dataSize
	}

	unimportantcan := newUnimportantCAN(ucPeriod, ucDeadline)
	unimportantcan.DataSize = dataSize

	return unimportantcan
}

func randomUnimportantCAN() (int, int) {
	// Use configured parameters if available, otherwise use defaults
	periods := unimportantCANParams.periods
	deadlines := unimportantCANParams.deadlines

	if len(periods) == 0 {
		periods = []int{50000, 100000, 150000}
	}
	if len(deadlines) == 0 {
		deadlines = []int{10000, 12000, 14000, 16000, 18000, 20000}
	}

	periodIdx := rng.IntN(len(periods))
	deadlineIdx := rng.IntN(len(deadlines))

	return periods[periodIdx], deadlines[deadlineIdx]
}

func randomCANDevicesForPath(canNodeSet []int) (int, int) {
	sourceIndex := rng.IntN(len(canNodeSet))
	sourceNode := canNodeSet[sourceIndex]

	// Create temporary set without source node
	tempSet := make([]int, 0, len(canNodeSet)-1)
	tempSet = append(tempSet, canNodeSet[:sourceIndex]...)
	tempSet = append(tempSet, canNodeSet[sourceIndex+1:]...)

	destIndex := rng.IntN(len(tempSet))
	destinationNode := tempSet[destIndex]

	return sourceNode - 2000, destinationNode - 1000 // source id: 1000+id  destination id: 2000+id
}
