package routes

import (
	"fmt"
	"src/network"
)

func Get_OSACO_Routing(network *network.Network, K int) *KTrees_set {
	ktrees_set := &KTrees_set{}
	v2v := &V2V{}

	for nth, flow := range network.Flow_Set.TSNFlows {
		tree := SteninerTree(v2v, network.Graph_Set.TSNGraphs[nth], flow.Source, flow.Destinations, network.Bytes_Rate)
		Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, network.Bytes_Rate)
		ktrees_set.TSNTrees = append(ktrees_set.TSNTrees, Ktrees)
	}
	fmt.Printf("Finish %d TSN streams routes\n", len(ktrees_set.TSNTrees))

	for nth, flow := range network.Flow_Set.AVBFlows {
		tree := SteninerTree(v2v, network.Graph_Set.AVBGraphs[nth], flow.Source, flow.Destinations, network.Bytes_Rate)
		Ktrees := KSpanningTree(v2v, tree, K, flow.Source, flow.Destinations, network.Bytes_Rate)
		ktrees_set.AVBTrees = append(ktrees_set.AVBTrees, Ktrees)
	}
	fmt.Printf("Finish %d AVB streams routes\n", len(ktrees_set.AVBTrees))

	return ktrees_set
}

func (ktrees_set *KTrees_set) Input_ktree_set() *KTrees_set {
	Input_ktree_set := &KTrees_set{}

	var (
		input_tsn_end int = len(ktrees_set.TSNTrees) / 2
		input_avb_end int = len(ktrees_set.AVBTrees) / 2
	)

	Input_ktree_set.TSNTrees = append(Input_ktree_set.TSNTrees, ktrees_set.TSNTrees[:input_tsn_end]...)
	Input_ktree_set.AVBTrees = append(Input_ktree_set.AVBTrees, ktrees_set.AVBTrees[:input_avb_end]...)

	return Input_ktree_set
}

func (ktrees_set *KTrees_set) BG_ktree_set() *KTrees_set {
	BG_ktree_set := &KTrees_set{}

	var (
		bg_tsn_start int = len(ktrees_set.TSNTrees) / 2
		bg_avb_start int = len(ktrees_set.AVBTrees) / 2
	)

	BG_ktree_set.TSNTrees = append(BG_ktree_set.TSNTrees, ktrees_set.TSNTrees[bg_tsn_start:]...)
	BG_ktree_set.AVBTrees = append(BG_ktree_set.AVBTrees, ktrees_set.AVBTrees[bg_avb_start:]...)

	return BG_ktree_set
}

func (ktrees_set *KTrees_set) Input_SteinerTree_set() *Trees_set {
	Input_tree_set := &Trees_set{}

	var (
		input_tsn_end int = len(ktrees_set.TSNTrees) / 2
		input_avb_end int = len(ktrees_set.AVBTrees) / 2
	)

	for _, ktree := range ktrees_set.TSNTrees[:input_tsn_end] {
		Input_tree_set.TSNTrees = append(Input_tree_set.TSNTrees, ktree.Trees[0])
	}
	for _, ktree := range ktrees_set.AVBTrees[:input_avb_end] {
		Input_tree_set.AVBTrees = append(Input_tree_set.AVBTrees, ktree.Trees[0])
	}

	return Input_tree_set
}

func (ktrees_set *KTrees_set) BG_SteinerTree_set() *Trees_set {
	BG_tree_set := &Trees_set{}

	var (
		bg_tsn_start int = len(ktrees_set.TSNTrees) / 2
		bg_avb_start int = len(ktrees_set.AVBTrees) / 2
	)

	for _, ktree := range ktrees_set.TSNTrees[bg_tsn_start:] {
		BG_tree_set.TSNTrees = append(BG_tree_set.TSNTrees, ktree.Trees[0])
	}
	for _, ktree := range ktrees_set.AVBTrees[bg_avb_start:] {
		BG_tree_set.AVBTrees = append(BG_tree_set.AVBTrees, ktree.Trees[0])
	}

	return BG_tree_set
}
