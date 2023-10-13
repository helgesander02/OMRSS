package routes

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

func (tree *Tree) GetNodeByID(id int) *Node {
	for _, node := range tree.Nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}
