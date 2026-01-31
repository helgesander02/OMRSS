package can

import (
	"src/internal/random"
)

// Global RNG instance (will be set by network package)
var rng *random.Generator

// SetRNG sets the random number generator for this package
func SetRNG(r *random.Generator) {
	rng = r
}

func config_ImportantCAN_Stream() *importantCAN {
	importantcan := new_importantCAN()

	return importantcan
}

func config_UnimportantCAN_Stream() *unimportantCAN {
	ucPeriod, ucDeadline := random_UnimportantCAN()
	unimportantcan := new_unimportantCAN(ucPeriod, ucDeadline)

	return unimportantcan
}

func random_UnimportantCAN() (int, int) {
	unimportantCANPeriodArr := []int{50000, 100000, 150000}
	unimportantCANDeadlineArr := []int{10000, 12000, 14000, 16000, 18000, 20000}

	periodIdx := rng.IntN(len(unimportantCANPeriodArr))
	deadlineIdx := rng.IntN(len(unimportantCANDeadlineArr))

	return unimportantCANPeriodArr[periodIdx], unimportantCANDeadlineArr[deadlineIdx]
}

func random_CAN_Devices_For_Path(canNodeSet []int) (int, int) {
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
