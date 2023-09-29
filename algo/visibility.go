package algo

import (
	"src/flow"
	"src/routes"
	"time"
)

func CompVB(ktrees_set *routes.KTrees_set, Input_flow_set *flow.Flows, BG_flow_set *flow.Flows, bytes_rate float64) *Visibility {
	var (
		preference   float64 = 2.
		bg_tsn_start int     = len(ktrees_set.TSNTrees) / 2
		bg_avb_start int     = len(ktrees_set.AVBTrees) / 2
	)

	visibility := &Visibility{}
	for nth, tsn_ktree := range ktrees_set.TSNTrees {
		var v []float64
		for kth, _ := range tsn_ktree.Trees {
			mult := 1.
			if nth >= bg_tsn_start && kth == 0 {
				mult = preference
			}

			value := mult / float64(len(tsn_ktree.Trees))
			v = append(v, value)
		}
		visibility.TSN_VB = append(visibility.TSN_VB, v)
	}

	for nth, avb_ktree := range ktrees_set.AVBTrees {
		var v []float64
		for kth, tree := range avb_ktree.Trees {
			mult := 1.
			if nth >= bg_avb_start && kth == 0 {
				mult = preference
			}

			if nth < bg_avb_start {
				value := mult / WCD(tree, Input_flow_set.AVBFlows[nth], bytes_rate)
				v = append(v, value)

			} else {
				value := mult / WCD(tree, BG_flow_set.AVBFlows[nth-bg_avb_start], bytes_rate)
				v = append(v, value)
			}
		}
		visibility.AVB_VB = append(visibility.AVB_VB, v)
	}

	return visibility
}

func WCD(tree *routes.Tree, flow *flow.Flow, bytes_rate float64) float64 {
	//path := Route_MaxPath(tree)
	per_hop := time.Duration(0)

	per_hop += transmit_avb_itself(flow.Streams[0].DataSize, bytes_rate)
	per_hop += interfere_from_be(bytes_rate)
	per_hop += interfere_from_avb()
	per_hop += interfere_from_tsn()

	return float64(per_hop)
}

func Route_MaxPath(tree *routes.Tree) []int {

	return nil
}

func transmit_avb_itself(datasize float64, bytes_rate float64) time.Duration {
	/// Maximum proportion of bandwidth that AVB streams can occupy.
	const MAX_AVB_SETTING float64 = 0.75
	nanoseconds := datasize * bytes_rate * MAX_AVB_SETTING
	duration := time.Duration(int64(nanoseconds))
	return duration
}

func interfere_from_be(bytes_rate float64) time.Duration {
	// Maximum number of bytes in a frame.
	const MTU float64 = 1500.
	nanoseconds := MTU * bytes_rate
	duration := time.Duration(int64(nanoseconds))
	return duration
}

// Don Pannell, "AVB Latency Math"
func interfere_from_avb() time.Duration {

	return time.Duration(0)
}

// Sune Mølgaard Laursen, Paul Pop, Wilfried Steiner, "Routing Optimization of AVB Streams in TSN Networks"
func interfere_from_tsn() time.Duration {

	return time.Duration(0)
}
