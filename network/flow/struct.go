package flow

import (
	"src/network/flow/can"
	"src/network/flow/tt"
)

type FlowSet struct {
	TSNFlows          []*tt.Flow
	AVBFlows          []*tt.Flow
	EncapsulateMethod []*can.Method
}

func newFlowSet() *FlowSet {
	return &FlowSet{}
}
