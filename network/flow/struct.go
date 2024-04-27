package flow

type TSN struct {
	Period   int     // 100~2000us up 500us
	Deadline int     // Period = Deadline
	DataSize float64 // 30~100bytes up 10bytes
}

func new_TSN(t_period int, t_datasize float64) *TSN {
	return &TSN{
		Period:   t_period,
		Deadline: t_period,
		DataSize: t_datasize,
	}
}

type AVB struct {
	Period   int     // 125us
	Deadline int     // 2000us
	DataSize float64 // 1000~1500bytes  up 100bytes
}

func new_AVB(a_datasize float64) *AVB {
	return &AVB{
		Period:   125,
		Deadline: 2000,
		DataSize: a_datasize,
	}
}

type Stream struct {
	Name        string
	ArrivalTime int
	DataSize    float64
	Deadline    int
	FinishTime  int
}

func new_Stream(name string, arrivalTime int, datasize float64, deadline int, finishTime int) *Stream {
	return &Stream{
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
	Streams      []*Stream
}

func new_Flow(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	return &Flow{
		Period:      period,
		Deadline:    deadline,
		DataSize:    datasize,
		HyperPeriod: HyperPeriod,
	}
}

type Flows struct {
	TSNFlows []*Flow
	AVBFlows []*Flow
}

func new_Flows() *Flows {
	return &Flows{}
}
