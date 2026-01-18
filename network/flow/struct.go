package flow

import (
	"src/network/flow/can"
	"src/network/flow/tt"
)

type Flow_Set struct {
	TSNFlows          []*tt.Flow
	AVBFlows          []*tt.Flow
	EncapsulateMethod []*can.Method
}

func new_Flow_Set() *Flow_Set {
	return &Flow_Set{}
}
