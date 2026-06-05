package network

import (
	"src/network/flow/can"
	"src/network/flow/tt"
	"src/network/topology"
	"src/pkg/config"
	"src/pkg/random"
)

func FillRNG(r *random.Generator) {
	tt.SetRNG(r)
	can.SetRNG(r)
	topology.FillRNG(r)
}

func FillTTParams(cfg *config.Config) {
	tt.SetTSNParams(
		cfg.Network.Flows.TSN.Params.Periods,
		cfg.Network.Flows.TSN.Params.DataSizes,
	)
	tt.SetAVBParams(
		cfg.Network.Flows.AVB.Params.Period,
		cfg.Network.Flows.AVB.Params.Deadline,
		cfg.Network.Flows.AVB.Params.DataSizes,
	)
}

func FillCANParams(cfg *config.Config) {
	can.SetImportantCANParams(
		cfg.Network.Flows.CAN.Params.Important.Period,
		cfg.Network.Flows.CAN.Params.Important.Deadline,
		cfg.Network.Flows.CAN.Params.Important.DataSize,
	)
	can.SetUnimportantCANParams(
		cfg.Network.Flows.CAN.Params.Unimportant.Periods,
		cfg.Network.Flows.CAN.Params.Unimportant.Deadlines,
		cfg.Network.Flows.CAN.Params.Unimportant.DataSize,
	)
}
