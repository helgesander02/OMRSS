package flow

type TSN struct {
	Period   int     // 100~2000us up500us
	Deadline int     // Period = Deadline
	DataSize float64 // 30~100bytes up10bytes
}

type AVB struct {
	Period   int     // 125us
	Deadline int     // 2000us
	DataSize float64 // 1000~1500bytes  up100bytes
}

type Stream struct {
	Name        string
	ArrivalTime int
	DataSize    float64
	Deadline    int
	FinishTime  int
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

type Flows struct {
	TSNFlows []*Flow
	AVBFlows []*Flow
}
