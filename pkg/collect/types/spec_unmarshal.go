package types

import "encoding/json"

func (d *Doc) UnmarshalJSON(data []byte) error {
	var deprecatedCollectors struct {
		Specs []Spec `json:"specs"`
	}
	if err := json.Unmarshal(data, &deprecatedCollectors); err != nil {
		return err
	}

	var deprecatedLifecycle struct {
		Lifecycle []LifecycleTask `json:"lifecycle"`
	}
	// ignore errors on this one, bceause lifecycle might be an array or an object,
	// it's not an error if unmarshalling fails
	json.Unmarshal(data, &deprecatedLifecycle)

	var currentCollectors struct {
		Collect Collect `json:"collect"`
	}
	if err := json.Unmarshal(data, &currentCollectors); err != nil {
		return err
	}

	var currentLifecycle struct {
		Lifecycle Lifecycle `json:"lifecycle"`
	}
	// Also not an error becaUSe it might have been the older, unversioned
	// array format above
	json.Unmarshal(data, &currentLifecycle)

	d.Collect.V1 = append(currentCollectors.Collect.V1, deprecatedCollectors.Specs...)
	d.Lifecycle = append(currentLifecycle.Lifecycle.V1, deprecatedLifecycle.Lifecycle...)
	return nil
}
