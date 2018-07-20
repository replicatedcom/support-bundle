package types

import "encoding/json"

func (d *Doc) UnmarshalJSON(data []byte) error {
	var intermediate struct {
		Collect   Collect         `json:"collect,omitempty"`
		Lifecycle []LifecycleTask `json:"lifecycle,omitempty"`
		Specs     []Spec          `json:"specs,omitempty"`
	}
	if err := json.Unmarshal(data, &intermediate); err != nil {
		return err
	}
	d.Collect.V1 = append(intermediate.Collect.V1, intermediate.Specs...)
	d.Lifecycle = intermediate.Lifecycle
	return nil
}
