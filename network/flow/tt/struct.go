package tt

type TSN struct {
	Period   int     // 100~2000us up 500us
	Deadline int     // Period = Deadline
	DataSize float64 // 30~100bytes up 10bytes
}

func newTSN(tPeriod int, tDatasize float64) *TSN {
	return &TSN{
		Period:   tPeriod,
		Deadline: tPeriod,
		DataSize: tDatasize,
	}
}

type AVB struct {
	Period   int     // 125us
	Deadline int     // 2000us
	DataSize float64 // 1000~1500bytes  up 100bytes
}

func newAVB(aDatasize float64) *AVB {
	// Use config values if set, otherwise use defaults
	period := avbPeriod
	if period == 0 {
		period = 125
	}

	deadline := avbDeadline
	if deadline == 0 {
		deadline = 2000
	}

	return &AVB{
		Period:   period,
		Deadline: deadline,
		DataSize: aDatasize,
	}
}

type Frame struct {
	Name        string
	ArrivalTime int
	DataSize    float64
	Deadline    int
	FinishTime  int
}

func newTTFrame(name string, arrivalTime int, datasize float64, deadline int, finishTime int) *Frame {
	return &Frame{
		Name:        name,
		ArrivalTime: arrivalTime,
		DataSize:    datasize,
		Deadline:    deadline,
		FinishTime:  finishTime,
	}
}

type Flow struct {
	Period       int
	Deadline     int
	DataSize     float64
	HyperPeriod  int
	Source       int
	Destinations []int
	Frames       []*Frame
}

func newTTFlow(period int, deadline int, datasize float64, hyperPeriod int) *Flow {
	return &Flow{
		Period:      period,
		Deadline:    deadline,
		DataSize:    datasize,
		HyperPeriod: hyperPeriod,
	}
}

func newTSNFlows() []*Flow {
	var tsnFlows []*Flow
	return tsnFlows
}

func newAVBFlows() []*Flow {
	var avbFlows []*Flow
	return avbFlows
}
