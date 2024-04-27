package routes

import (
	"fmt"
	"src/network"
)

var v2v *V2V = &V2V{}

func Get_SteninerTree_Routing(network *network.Network) *Trees_set {
	Trees_set := new_Trees_Set()

	for nth, flow := range network.Flow_Set.TSNFlows {
		tree := SteninerTree(v2v, network.Graph_Set.TSNGraphs[nth], flow.Source, flow.Destinations, network.BytesRate)
		Trees_set.TSNTrees = append(Trees_set.TSNTrees, tree)
	}
	fmt.Printf("Finish Steniner Tree %d TSN streams routing\n", len(Trees_set.TSNTrees))

	for nth, flow := range network.Flow_Set.AVBFlows {
		tree := SteninerTree(v2v, network.Graph_Set.AVBGraphs[nth], flow.Source, flow.Destinations, network.BytesRate)
		Trees_set.AVBTrees = append(Trees_set.AVBTrees, tree)
	}
	fmt.Printf("Finish Steniner Tree %d AVB streams routing\n", len(Trees_set.AVBTrees))

	return Trees_set
}

func Get_DistanceTree_Routing(network *network.Network) *Trees_set {
	Trees_set := new_Trees_Set()

	for nth, flow := range network.Flow_Set.TSNFlows {
		tree := DistanceTree(network.Graph_Set.TSNGraphs[nth], flow.Source, flow.Destinations, network.BytesRate)
		Trees_set.TSNTrees = append(Trees_set.TSNTrees, tree)
	}
	fmt.Printf("Finish Distance Tree %d TSN streams routing\n", len(Trees_set.TSNTrees))

	for nth, flow := range network.Flow_Set.AVBFlows {
		tree := DistanceTree(network.Graph_Set.AVBGraphs[nth], flow.Source, flow.Destinations, network.BytesRate)
		Trees_set.AVBTrees = append(Trees_set.AVBTrees, tree)
	}
	fmt.Printf("Finish Distance Tree %d AVB streams routing\n", len(Trees_set.AVBTrees))

	return Trees_set
}

func (trees_set *Trees_set) Input_Tree_set() *Trees_set {
	Input_tree_set := new_Trees_Set()

	var (
		input_tsn_end int = len(trees_set.TSNTrees) / 2
		input_avb_end int = len(trees_set.AVBTrees) / 2
	)

	Input_tree_set.TSNTrees = append(Input_tree_set.TSNTrees, trees_set.TSNTrees[:input_tsn_end]...)
	Input_tree_set.AVBTrees = append(Input_tree_set.AVBTrees, trees_set.AVBTrees[:input_avb_end]...)

	return Input_tree_set
}

func (trees_set *Trees_set) BG_Tree_set() *Trees_set {
	BG_tree_set := new_Trees_Set()

	var (
		bg_tsn_start int = len(trees_set.TSNTrees) / 2
		bg_avb_start int = len(trees_set.AVBTrees) / 2
	)
	BG_tree_set.TSNTrees = append(BG_tree_set.TSNTrees, trees_set.TSNTrees[bg_tsn_start:]...)
	BG_tree_set.AVBTrees = append(BG_tree_set.AVBTrees, trees_set.AVBTrees[bg_avb_start:]...)

	return BG_tree_set
}

func Get_OSACO_Routing(network *network.Network, SMT *Trees_set, K int) *KTrees_set {
	ktrees_set := new_KTrees_Set()

	for nth, flow := range network.Flow_Set.TSNFlows {
		Ktrees := KSpanningTree(v2v, SMT.TSNTrees[nth], K, flow.Source, flow.Destinations, network.BytesRate)
		ktrees_set.TSNTrees = append(ktrees_set.TSNTrees, Ktrees)
	}
	fmt.Printf("Finish OSACO %d TSN streams routing\n", len(ktrees_set.TSNTrees))

	for nth, flow := range network.Flow_Set.AVBFlows {
		Ktrees := KSpanningTree(v2v, SMT.AVBTrees[nth], K, flow.Source, flow.Destinations, network.BytesRate)
		ktrees_set.AVBTrees = append(ktrees_set.AVBTrees, Ktrees)
	}
	fmt.Printf("Finish OSACO %d AVB streams routing\n", len(ktrees_set.AVBTrees))

	return ktrees_set
}

func (ktrees_set *KTrees_set) Input_ktree_set() *KTrees_set {
	Input_ktree_set := new_KTrees_Set()

	var (
		input_tsn_end int = len(ktrees_set.TSNTrees) / 2
		input_avb_end int = len(ktrees_set.AVBTrees) / 2
	)

	Input_ktree_set.TSNTrees = append(Input_ktree_set.TSNTrees, ktrees_set.TSNTrees[:input_tsn_end]...)
	Input_ktree_set.AVBTrees = append(Input_ktree_set.AVBTrees, ktrees_set.AVBTrees[:input_avb_end]...)

	return Input_ktree_set
}

func (ktrees_set *KTrees_set) BG_ktree_set() *KTrees_set {
	BG_ktree_set := new_KTrees_Set()

	var (
		bg_tsn_start int = len(ktrees_set.TSNTrees) / 2
		bg_avb_start int = len(ktrees_set.AVBTrees) / 2
	)

	BG_ktree_set.TSNTrees = append(BG_ktree_set.TSNTrees, ktrees_set.TSNTrees[bg_tsn_start:]...)
	BG_ktree_set.AVBTrees = append(BG_ktree_set.AVBTrees, ktrees_set.AVBTrees[bg_avb_start:]...)

	return BG_ktree_set
}
