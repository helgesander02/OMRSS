package can

import (
	"time"
)

type importantCAN struct {
	Period   int
	Deadline int
	DataSize float64
}

func newImportantCANWithParams() *importantCAN {
	return &importantCAN{
		Period:   importantcanPeriods,   // 5000us
		Deadline: importantcanDeadlines, // Period = Deadline
		DataSize: importantcanDataSizes, // 16bytes
	}
}

type unimportantCAN struct {
	Period   int
	Deadline int
	DataSize float64
}

func newUnimportantCAN(period int, deadline int) *unimportantCAN {
	return &unimportantCAN{
		Period:   period,                  // 50000~150000us up 50000us
		Deadline: deadline,                // 10000~20000us up 2000us
		DataSize: unimportantcanDataSizes, // 16bytes
	}
}

type Frame struct {
	Name        string
	ArrivalTime int
	DataSize    float64
	Deadline    int
	FinishTime  int
}

func newCANFrame(name string, arrivalTime int, datasize float64, deadline int, finishTime int) *Frame {
	return &Frame{
		Name:        name,
		ArrivalTime: arrivalTime,
		DataSize:    datasize,
		Deadline:    deadline,
		FinishTime:  finishTime,
	}
}

func createCAN2TTFrame(arrivalTime int, deadline int, datasize float64) *Frame {
	newFrame := &Frame{
		ArrivalTime: arrivalTime,
		Deadline:    deadline,
		DataSize:    datasize,
		FinishTime:  arrivalTime + deadline,
	}

	return newFrame
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

func newCANFlow(period int, deadline int, datasize float64, hyperPeriod int) *Flow {
	return &Flow{
		Period:      period,
		Deadline:    deadline,
		DataSize:    datasize,
		HyperPeriod: hyperPeriod,
	}
}

func newCAN2TTFlow() *Flow {
	return &Flow{}
}

type Method struct {
	MethodName    string
	CAN2TTFlows   []*Flow
	BytesSent     float64
	TTFrameCount  int
	CAN2TTO1Drop  int
	CANAreaO1Drop int
	CAN2TTDelay   time.Duration
}

func newMethod(methodName string) *Method {
	return &Method{
		MethodName: methodName,
	}
}

func newMethodSet() []*Method {
	var methodSet []*Method
	return methodSet
}
