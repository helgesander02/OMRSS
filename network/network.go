package network

type Networks interface {
	GenerateNetwork()
	ShowNetwork()
}

func NewNetworks(topologyName string, bgTSN int, bgAVB int, inputTSN int, inputAVB int, importantCAN int, unimportantCAN int, hyperperiod int, bandwidth float64) map[string]Networks {
	// network1 ...
	OMACO := new_OMACO_Network(topologyName, bgTSN, bgAVB, inputTSN, inputAVB, hyperperiod, bandwidth)
	OSRO := new_OSRO_Network(topologyName, bgTSN, bgAVB, inputTSN, inputAVB, importantCAN, unimportantCAN, hyperperiod, bandwidth)

	// Look-up table method
	Networks := map[string]Networks{
		"omaco": OMACO,
		"osro":  OSRO,
		//network3,
		// ...
	}

	return Networks
}
