package topology

import "encoding/json"

// DeepCopy Topology
func (t1 *Topology) TopologyDeepCopy() *Topology {
	if buf, err := json.Marshal(t1); err != nil {
		return nil
	} else {
		t2 := &Topology{}
		if err = json.Unmarshal(buf, t2); err != nil {
			return nil
		}
		return t2
	}
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
