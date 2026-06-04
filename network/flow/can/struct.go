package can

import (
	"time"
)

type importantCAN struct {
	Period   int
	Deadline int
	DataSize float64
}

func newImportantCAN() *importantCAN {
	return &importantCAN{
		Period:   5000, // 5000us (default)
		Deadline: 5000, // Period = Deadline (default)
		DataSize: 16,   // 16bytes (default)
	}
}

func newImportantCANWithParams(period int, deadline int, dataSize float64) *importantCAN {
	return &importantCAN{
		Period:   period,
		Deadline: deadline,
		DataSize: dataSize,
	}
}

type unimportantCAN struct {
	Period   int
	Deadline int
	DataSize float64
}

func newUnimportantCAN(ucPeriod int, ucDeadline int) *unimportantCAN {
	return &unimportantCAN{
		Period:   ucPeriod,   // 50000~150000us up 50000us
		Deadline: ucDeadline, // 10000~20000us up 2000us
		DataSize: 16,         // 16bytes
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
	Period      int
	Deadline    int
	DataSize    float64
	HyperPeriod int
	Source      int
	Destination int
	Frames      []*Frame
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
	CAN2TSNDelay  time.Duration
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
