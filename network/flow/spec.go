package flow

import (
	"crypto/rand"
	"math/big"
)

func config_TSN_Stream() *TSN {
	t_period, t_datasize := random_TSN()
	tsn := new_TSN(t_period, t_datasize)

	return tsn
}

func config_AVB_Stream() *AVB {
	a_datasize := random_AVB()
	avb := new_AVB(a_datasize)

	return avb
}

func random_TSN() (int, float64) {
	tsn_period_arr := []int{100, 500, 1000, 1500, 2000}
	tsn_datasize_arr := []float64{30., 40., 50., 60., 70., 80., 90., 100.}
	period_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(tsn_period_arr))))
	datasize_rng, _ := rand.Int(rand.Reader, big.NewInt(int64((len(tsn_datasize_arr)))))

	return tsn_period_arr[period_rng.Int64()], tsn_datasize_arr[datasize_rng.Int64()]
}

func random_AVB() float64 {
	avb_datasize_arr := []float64{1000., 1100., 1200., 1300., 1400., 1500.}
	datasize_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(avb_datasize_arr))))

	return avb_datasize_arr[datasize_rng.Int64()]
}

func random_TT_Devices_For_Tree(Nnode int) (int, []int) {
	// Talker
	source, _ := rand.Int(rand.Reader, big.NewInt(int64(Nnode)))

	// Listener
	destinations := []int{}
	for i := 0; i < Nnode; i++ {
		if i != int(source.Int64()) {
			destinations = append(destinations, i+2000)
		}
	}

	numDestinations, _ := rand.Int(rand.Reader, big.NewInt(2))
	max := big.NewInt(int64(Nnode - 1)) // 10 (0~9) - source = 9
	num, _ := rand.Int(rand.Reader, max.Sub(max, big.NewInt(3)))
	n := num.Add(num, big.NewInt(3)).Int64()
	numDestinations = numDestinations.Add(numDestinations, big.NewInt(n-1))

	selectedDestinations := []int{}
	for i := 0; i < int(numDestinations.Int64()); i++ {
		// Randomly selects an element from the 'destinations' slice.
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(destinations))))
		selectedIndex := int(randIndex.Int64())
		selectedDestinations = append(selectedDestinations, destinations[selectedIndex])
		// To prevent repeated selection, remove the selected element from the 'destinations' slice.
		destinations = append(destinations[:selectedIndex], destinations[selectedIndex+1:]...)
	}

	return int(source.Int64()) + 1000, selectedDestinations
}
