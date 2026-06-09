package flow

import (
	"fmt"

	"src/network/flow/can"
	"src/pkg/logger"
)

func (flows *FlowSet) ShowFrame() {
	TSNFlows := flows.TSNFlows
	AVBFlows := flows.AVBFlows
	number := 1
	for _, flow := range TSNFlows {
		name := fmt.Sprint("TSNflow", number)
		logger.Println(name)
		for _, frame := range flow.Frames {
			logger.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				frame.Name, frame.ArrivalTime, frame.DataSize, frame.Deadline, frame.FinishTime)
		}
		number += 1

		break
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBflow", number)
		logger.Println(name)
		for _, frame := range flow.Frames {
			logger.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				frame.Name, frame.ArrivalTime, frame.DataSize, frame.Deadline, frame.FinishTime)
		}
		number += 1

		break
	}
}

func (flows *FlowSet) ShowTTFlow() {
	TSNFlows := flows.TSNFlows
	AVBFlows := flows.AVBFlows
	number := 1
	for _, flow := range TSNFlows {
		name := fmt.Sprint("TTflow", number)
		logger.Printf("Source: %d\n", flow.Source)
		logger.Printf("Destinations: %v\n", flow.Destinations)
		logger.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBflow", number)
		logger.Printf("Source: %d\n", flow.Source)
		logger.Printf("Destinations: %v\n", flow.Destinations)
		logger.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}
}

// ShowCANFlow dumps every encapsulation method's CAN2TT flow set.
// Summary stats and the cross-method comparison table always print.
// When verbose is true, also dumps per-flow params and per-emit TT frame details.
func (flows *FlowSet) ShowCANFlow(verbose bool) {
	if len(flows.EncapsulateMethod) == 0 {
		return
	}

	logger.Println()
	logger.Println("=== CAN2TT Encapsulation ===")
	for _, method := range flows.EncapsulateMethod {
		logger.Println()
		logger.Printf("Method: %s\n", method.MethodName)
		// All Method fields, drops first because that's what matters for AR.
		logger.Printf("  CAN2TTO1Drop  : %d   (CAN frames overdue while waiting in gateway queue)\n", method.CAN2TTO1Drop)
		logger.Printf("  CANAreaO1Drop : %d   (CAN-side drops; not tracked yet, always 0)\n", method.CANAreaO1Drop)
		logger.Printf("  CAN2TTFlows   : %d\n", len(method.CAN2TTFlows))
		logger.Printf("  TTFrameCount  : %d\n", method.TTFrameCount)
		logger.Printf("  BytesSent     : %.0f B\n", method.BytesSent)
		logger.Printf("  CAN2TTDelay   : %v   (encap CPU time)\n", method.CAN2TTDelay)

		if !verbose {
			continue
		}
		for i, flow := range method.CAN2TTFlows {
			name := fmt.Sprint("CAN2TTflow", i+1)
			logger.Printf("  %s  src=%d → dst=%d  period=%d us  deadline=%d us  avg_payload=%.1f B  emits=%d\n",
				name, flow.Source, flow.Destination, flow.Period, flow.Deadline, flow.DataSize, len(flow.Frames))
			// Per-emit detail: payload = frame.DataSize - HeaderBytes (42).
			for j, frame := range flow.Frames {
				logger.Printf("      emit#%d  t=%d us  size=%.0f B (payload=%.0f)  deadline=%d  finish=%d\n",
					j+1, frame.ArrivalTime, frame.DataSize, frame.DataSize-can.HeaderBytes, frame.Deadline, frame.FinishTime)
			}
		}
	}

	// Summary comparison table: one row per method, easy to scan side by side.
	logger.Println()
	logger.Println("=== CAN2TT Method Comparison ===")
	logger.Printf("%-10s %-7s %-7s %-10s %-12s %-7s %-12s\n",
		"Method", "Flows", "Emits", "Bytes", "Delay", "Drops", "AvgPayload")
	logger.Printf("%-10s %-7s %-7s %-10s %-12s %-7s %-12s\n",
		"------", "-----", "-----", "-----", "-----", "-----", "----------")
	for _, method := range flows.EncapsulateMethod {
		// Avg payload per emit, excluding the 42-byte TT header.
		var avgPayload float64
		if method.TTFrameCount > 0 {
			avgPayload = method.BytesSent/float64(method.TTFrameCount) - can.HeaderBytes
		}
		logger.Printf("%-10s %-7d %-7d %-10.0f %-12v %-7d %-10.1f B\n",
			method.MethodName,
			len(method.CAN2TTFlows),
			method.TTFrameCount,
			method.BytesSent,
			method.CAN2TTDelay,
			method.CAN2TTO1Drop,
			avgPayload,
		)
	}
}

func (flows *FlowSet) ShowFlows() {
	// Display all flows.
	logger.Printf("Total Flows:%d ( TSN Flows:%d  AVB Flows:%d )\n",
		len(flows.TSNFlows)+len(flows.AVBFlows), len(flows.TSNFlows), len(flows.AVBFlows))
}
