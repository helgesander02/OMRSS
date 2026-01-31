package can

import (
	"time"
)

type importantCAN struct {
	Period   int
	Deadline int
	DataSize float64
}

func new_importantCAN() *importantCAN {
	return &importantCAN{
		Period:   5000, // 5000us
		Deadline: 5000, // Period = Deadline
		DataSize: 16,   // 16bytes
	}
}

type unimportantCAN struct {
	Period   int
	Deadline int
	DataSize float64
}

func new_unimportantCAN(ucPeriod int, ucDeadline int) *unimportantCAN {
	return &unimportantCAN{
		Period:   ucPeriod,   // 50000~150000us up 50000us
		Deadline: ucDeadline, // 10000~20000us up 2000us
		DataSize: 16,         // 16bytes
	}
}

type Stream struct {
	Name        string
	ArrivalTime int
	DataSize    float64
	Deadline    int
	FinishTime  int
}

func new_CANStream(name string, arrivalTime int, datasize float64, deadline int, finishTime int) *Stream {
	return &Stream{
		Name:        name,
		ArrivalTime: arrivalTime,
		DataSize:    datasize,
		Deadline:    deadline,
		FinishTime:  finishTime,
	}
}

func createCAN2TTStream(arrival_time int, deadline int, datasize float64) *Stream {
	newStream := &Stream{
		ArrivalTime: arrival_time,
		Deadline:    deadline,
		DataSize:    datasize,
		FinishTime:  arrival_time + deadline,
	}

	return newStream
}

type Flow struct {
	Period      int
	Deadline    int
	DataSize    float64
	HyperPeriod int
	Source      int
	Destination int
	Streams     []*Stream
}

func new_CANFlow(period int, deadline int, datasize float64, hyperPeriod int) *Flow {
	return &Flow{
		Period:      period,
		Deadline:    deadline,
		DataSize:    datasize,
		HyperPeriod: hyperPeriod,
	}
}

func new_CAN2TTFlow() *Flow {
	return &Flow{}
}

type Method struct {
	Method_Name      string
	CAN2TTFlows      []*Flow
	BytesSent        float64
	TTFrameCount     int
	CAN2TT_O1_Drop   int
	CAN_Area_O1_Drop int
	CAN2TSN_Delay    time.Duration
}

func new_Method(method_name string) *Method {
	return &Method{
		Method_Name: method_name,
	}
}

func new_Method_Set() []*Method {
	var method_set []*Method
	return method_set
}
