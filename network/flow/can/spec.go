package can

import (
	"src/pkg/random"
)

// random variables
var rng *random.Generator

func SetRNG(r *random.Generator) {
	rng = r
}

// parameters from config
var (
	importantcanPeriods   int
	importantcanDeadlines int
	importantcanDataSizes float64

	unimportantcanPeriods   []int
	unimportantcanDeadlines []int
	unimportantcanDataSizes float64
)

func SetImportantCANParams(period int, deadline int, dataSize float64) {
	importantcanPeriods = period
	importantcanDeadlines = deadline
	importantcanDataSizes = dataSize
}

func SetUnimportantCANParams(periods []int, deadlines []int, dataSize float64) {
	unimportantcanPeriods = periods
	unimportantcanDeadlines = deadlines
	unimportantcanDataSizes = dataSize
}

func configImportantCANFrame() *importantCAN {
	return newImportantCANWithParams()
}

func configUnimportantCANFrame() *unimportantCAN {
	ucPeriod, ucDeadline := randomUnimportantCAN()
	unimportantcan := newUnimportantCAN(ucPeriod, ucDeadline)

	return unimportantcan
}

func randomUnimportantCAN() (int, int) {
	periodIdx := rng.IntN(len(unimportantcanPeriods))
	deadlineIdx := rng.IntN(len(unimportantcanDeadlines))

	return unimportantcanPeriods[periodIdx], unimportantcanDeadlines[deadlineIdx]
}

func randomCANDevices(canNodeSet []int) (int, int) {
	sourceIndex := rng.IntN(len(canNodeSet))
	sourceNode := canNodeSet[sourceIndex]

	tempSet := make([]int, 0, len(canNodeSet)-1)
	tempSet = append(tempSet, canNodeSet[:sourceIndex]...)
	tempSet = append(tempSet, canNodeSet[sourceIndex+1:]...)

	destIndex := rng.IntN(len(tempSet))
	destinationNode := tempSet[destIndex]

	return sourceNode - 2000, destinationNode - 1000 // source id: 1000+id  destination id: 2000+id
}
