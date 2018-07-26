package types

import "encoding/json"

func (d *Doc) UnmarshalJSON(data []byte) error {
	var intermediate struct {
		Collect   Collect         `json:"collect"`
		Lifecycle []LifecycleTask `json:"lifecycle"`
		Specs     []Spec          `json:"specs"`
	}
	if err := json.Unmarshal(data, &intermediate); err != nil {
		return err
	}
	d.Collect.V1 = append(intermediate.Collect.V1, intermediate.Specs...)
	d.Lifecycle = intermediate.Lifecycle
	return nil
}
