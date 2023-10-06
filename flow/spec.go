package flow

import (
	"crypto/rand"
	"math/big"
)

func TSN_stream() *TSN {
	t_period, t_datasize := tsn_rondam()

	tsn := &TSN{
		Period:   t_period,   // 100~2000us up500us
		Deadline: t_period,   // Period = Deadline
		DataSize: t_datasize, // 30~100bytes up10bytes
	}

	return tsn
}

func AVB_stream() *AVB {
	a_datasize := avb_rondam()

	avb := &AVB{
		Period:   125,        // 125us
		Deadline: 2000,       // 2000us
		DataSize: a_datasize, // 1000~1500bytes  up100bytes
	}

	return avb
}

func tsn_rondam() (int, float64) {
	tsn_period_arr := []int{100, 500, 1000, 1500, 2000}
	tsn_datasize_arr := []float64{30., 40., 50., 60., 70., 80., 90., 100.}

	period_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(tsn_period_arr))))
	datasize_rng, _ := rand.Int(rand.Reader, big.NewInt(int64((len(tsn_datasize_arr)))))
	return tsn_period_arr[period_rng.Int64()], tsn_datasize_arr[datasize_rng.Int64()]

}

func avb_rondam() float64 {
	avb_datasize_arr := []float64{1000., 1100., 1200., 1300., 1400., 1500.}

	datasize_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(avb_datasize_arr))))
	return avb_datasize_arr[datasize_rng.Int64()]
}

func Random_Devices(Nnode int) (int, []int) {
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
