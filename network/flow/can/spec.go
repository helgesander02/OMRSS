package can

import (
	"crypto/rand"
	"math/big"
)

func config_ImportantCAN_Stream() *importantCAN {
	importantcan := new_importantCAN()

	return importantcan
}

func config_UnimportantCAN_Stream() *unimportantCAN {
	uc_period, uc_deadline := random_UnimportantCAN()
	unimportantcan := new_unimportantCAN(uc_period, uc_deadline)

	return unimportantcan
}

func random_UnimportantCAN() (int, int) {
	unimportantCAN_period_arr := []int{50000, 100000, 150000}
	unimportantCAN_deadline_arr := []int{10000, 12000, 14000, 16000, 18000, 20000}
	period_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(unimportantCAN_period_arr))))
	deadline_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(unimportantCAN_deadline_arr))))

	return unimportantCAN_period_arr[period_rng.Int64()], unimportantCAN_deadline_arr[deadline_rng.Int64()]
}

func random_CAN_Devices_For_Path(CAN_Node_Set []int) (int, int) {
	sourceIndexBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(CAN_Node_Set))))
	sourceIndex := int(sourceIndexBig.Int64())
	sourceNode := CAN_Node_Set[sourceIndex]

	tempSet := make([]int, 0, len(CAN_Node_Set)-1)
	tempSet = append(tempSet, CAN_Node_Set[:sourceIndex]...)
	tempSet = append(tempSet, CAN_Node_Set[sourceIndex+1:]...)

	destIndexBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(tempSet))))
	destIndex := int(destIndexBig.Int64())
	destinationNode := tempSet[destIndex]

	return sourceNode - 2000, destinationNode - 1000 // source id: 1000+id  destination id: 2000+id
}
