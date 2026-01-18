package topology

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
)

// DeepCopy Topology
func (t1 *Topology) TopologyDeepCopy() *Topology {
	if buf, err := json.Marshal(t1); err != nil {
		return nil
	} else {
		t2 := new_Topology()
		if err = json.Unmarshal(buf, t2); err != nil {
			return nil
		}
		return t2
	}
}

func (t *Topology) Select_CAN_Node_Set() []int {
	const count = 5
	result := make([]int, 0, count)
	used := make(map[int]bool)

	for len(result) < count {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(t.Nodes))))
		if err != nil {
			return nil
		}

		index := int(randIndex.Int64())
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

func (t *Topology) RemoveEdge(u, v int) {
	n := t.GetNodeByID(u)
	if n == nil {
		return
	}
	for i, connection := range n.Connections {
		if connection.ToNodeID == v {
			n.Connections = append(n.Connections[:i], n.Connections[i+1:]...)
			return
		}
	}
}

func (t *Topology) GetListenerAndTalker(source int, destination int) bool {
	if t.Talker[0].ID == source && t.Listener[0].ID == destination {
		return true
	}
	return false
}

func (t *Topology) GetListenerAndTalkerSet(source int, destinations []int) bool {
	if t.Talker[0].ID == source && len(t.Listener) == len(destinations) {
		for _, desdestination := range destinations {
			if t.GetNodeByID(desdestination) == nil {
				return false
			}
		}
		return true
	}
	return false
}
