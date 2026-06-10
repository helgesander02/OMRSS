package routes

import (
	"src/network"
	"src/pkg/config"
	"src/pkg/logger"
)

var v2v *V2V = &V2V{} // v2v is all paths connecting multiple terminals to terminals.

func Get_SteninerTree_Routing(network *network.Network, cfg *config.Config) *TreesSet {
	TreesSet := newTreesSet()

	for _, flow := range network.FlowSet.TSNFlows {
		topo := network.GraphSet.Get(flow.Source, flow.Destinations)
		tree := SteninerTree(v2v, topo, flow.Source, flow.Destinations, cfg.Network.ByteRate)
		TreesSet.TSNTrees = append(TreesSet.TSNTrees, tree)
	}
	logger.Printf("Finish Steniner Tree %d TSN streams routing\n", len(TreesSet.TSNTrees))

	for _, flow := range network.FlowSet.AVBFlows {
		topo := network.GraphSet.Get(flow.Source, flow.Destinations)
		tree := SteninerTree(v2v, topo, flow.Source, flow.Destinations, cfg.Network.ByteRate)
		TreesSet.AVBTrees = append(TreesSet.AVBTrees, tree)
	}
	logger.Printf("Finish Steniner Tree %d AVB streams routing\n", len(TreesSet.AVBTrees))

	return TreesSet
}

func Get_DistanceTree_Routing(network *network.Network, cfg *config.Config) *TreesSet {
	TreesSet := newTreesSet()

	for _, flow := range network.FlowSet.TSNFlows {
		topo := network.GraphSet.Get(flow.Source, flow.Destinations)
		tree := DistanceTree(topo, flow.Source, flow.Destinations, cfg.Network.ByteRate)
		TreesSet.TSNTrees = append(TreesSet.TSNTrees, tree)
	}
	logger.Printf("Finish Distance Tree %d TSN streams routing\n", len(TreesSet.TSNTrees))

	for _, flow := range network.FlowSet.AVBFlows {
		topo := network.GraphSet.Get(flow.Source, flow.Destinations)
		tree := DistanceTree(topo, flow.Source, flow.Destinations, cfg.Network.ByteRate)
		TreesSet.AVBTrees = append(TreesSet.AVBTrees, tree)
	}
	logger.Printf("Finish Distance Tree %d AVB streams routing\n", len(TreesSet.AVBTrees))

	return TreesSet
}

func (trees_set *TreesSet) InputTreeSet(bg_tsn_end int, bg_avb_end int) *TreesSet {
	Input_tree_set := newTreesSet()

	Input_tree_set.TSNTrees = append(Input_tree_set.TSNTrees, trees_set.TSNTrees[bg_tsn_end:]...)
	Input_tree_set.AVBTrees = append(Input_tree_set.AVBTrees, trees_set.AVBTrees[bg_avb_end:]...)

	return Input_tree_set
}

func (trees_set *TreesSet) BGTreeSet(bg_tsn_end int, bg_avb_end int) *TreesSet {
	BG_tree_set := newTreesSet()

	BG_tree_set.TSNTrees = append(BG_tree_set.TSNTrees, trees_set.TSNTrees[:bg_tsn_end]...)
	BG_tree_set.AVBTrees = append(BG_tree_set.AVBTrees, trees_set.AVBTrees[:bg_avb_end]...)

	return BG_tree_set
}

func Get_OSACO_Routing(network *network.Network, cfg *config.Config, SMT *TreesSet, K int, Method_Number int) *KTreesSet {
	ktrees_set := newKTreesSet()

	for nth, flow := range network.FlowSet.TSNFlows {
		Ktrees := KSpanningTree(v2v, SMT.TSNTrees[nth], K, flow.Source, flow.Destinations, cfg.Network.ByteRate, Method_Number)
		ktrees_set.TSNTrees = append(ktrees_set.TSNTrees, Ktrees)
	}
	logger.Printf("Finish OSACO %d TSN streams routing\n", len(ktrees_set.TSNTrees))

	for nth, flow := range network.FlowSet.AVBFlows {
		Ktrees := KSpanningTree(v2v, SMT.AVBTrees[nth], K, flow.Source, flow.Destinations, cfg.Network.ByteRate, Method_Number)
		ktrees_set.AVBTrees = append(ktrees_set.AVBTrees, Ktrees)
	}
	logger.Printf("Finish OSACO %d AVB streams routing\n", len(ktrees_set.AVBTrees))

	return ktrees_set
}

func (ktrees_set *KTreesSet) Input_ktree_set(bg_tsn_end int, bg_avb_end int) *KTreesSet {
	Input_ktree_set := newKTreesSet()

	Input_ktree_set.TSNTrees = append(Input_ktree_set.TSNTrees, ktrees_set.TSNTrees[bg_tsn_end:]...)
	Input_ktree_set.AVBTrees = append(Input_ktree_set.AVBTrees, ktrees_set.AVBTrees[bg_tsn_end:]...)

	return Input_ktree_set
}

func (ktrees_set *KTreesSet) BG_ktree_set(bg_tsn_end int, bg_avb_end int) *KTreesSet {
	BG_ktree_set := newKTreesSet()

	BG_ktree_set.TSNTrees = append(BG_ktree_set.TSNTrees, ktrees_set.TSNTrees[:bg_tsn_end]...)
	BG_ktree_set.AVBTrees = append(BG_ktree_set.AVBTrees, ktrees_set.AVBTrees[:bg_tsn_end]...)

	return BG_ktree_set
}
