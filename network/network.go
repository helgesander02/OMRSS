package network

type Networks interface {
	Generate_Network()
	Show_Network()
}

func New_Networks(topology_name string, bg_tsn int, bg_avb int, input_tsn int, input_avb int, hyperperiod int, bandwidth float64) map[string]Networks {
	// network1 ...
	OMACO := new_OMACO_Network(topology_name, bg_tsn, bg_avb, input_tsn, input_avb, hyperperiod, bandwidth)

	// Look-up table method
	Networks := map[string]Networks{
		"omaco": OMACO,
		//network2,
		//network3,
		// ...
	}

	return Networks
}
