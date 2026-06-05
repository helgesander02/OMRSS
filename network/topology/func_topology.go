package topology

import (
	"encoding/json"
)

func (t1 *Topology) TopologyDeepCopy() *Topology {
	if buf, err := json.Marshal(t1); err != nil {
		return nil
	} else {
		t2 := NewTopology()
		if err = json.Unmarshal(buf, t2); err != nil {
			return nil
		}
		return t2
	}
}

func (t *Topology) SelectCANNodes(count int) []int {
	result := make([]int, 0, count)
	used := make(map[int]bool)

	for len(result) < count {
		index := rng.IntN(len(t.Nodes))
		node := t.Nodes[index]

		if !used[node.ID] {
			result = append(result, node.ID)
			used[node.ID] = true
		}
	}

	return result
}

func (t *Topology) GetNodeByID(id int) *Node {
	for _, node := range t.Talker {
		if node.ID == id {
			return node
		}
	}
	for _, node := range t.Switch {
		if node.ID == id {
			return node
		}
	}
	for _, node := range t.Listener {
		if node.ID == id {
			return node
		}
	}
	for _, node := range t.Nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}

func (t *Topology) GetListenerAndTalker(source int, destination int) bool {
	if t.Talker[0].ID == source && t.Listener[0].ID == destination {
		return true
	}
	return false
}

func (t *Topology) GetListenerAndTalkerSet(source int, destinations []int) bool {
	if t.Talker[0].ID == source && len(t.Listener) == len(destinations) {
		for _, desDestination := range destinations {
			if t.GetNodeByID(desDestination) == nil {
				return false
			}
		}
		return true
	}
	return false
}
