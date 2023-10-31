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
