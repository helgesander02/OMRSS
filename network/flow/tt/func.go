package tt

import "encoding/json"

func (f *Flow) deepcopyFlow() *Flow {
	if f == nil {
		return nil
	}

	data, err := json.Marshal(f)
	if err != nil {
		return nil
	}

	dst := &Flow{}
	if err := json.Unmarshal(data, dst); err != nil {
		return nil
	}
	return dst
}
