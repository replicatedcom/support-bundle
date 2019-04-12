package v1

import (
	"errors"
	"text/template"
)

var (
	SelinuxNotEnforcingRequirementSpecTemplate *template.Template
)

func init() {
	SelinuxNotEnforcingRequirementSpecTemplate = mustAsset(
		"SelinuxNotEnforcingRequirementSpec", "requirements/selinux_notenforcing.yml")
}

type SelinuxNotEnforcingRequirement struct {
}

func (a *SelinuxNotEnforcingRequirement) Validate(spec AnalyzerSpec) error {
	if len(spec.CollectRefs) == 0 {
		return errors.New("collect_refs required")
	}
	return nil
}

func (a *SelinuxNotEnforcingRequirement) GetRawSpec() (RawSpec, error) {
	return unmarshalSpec(SelinuxNotEnforcingRequirementSpecTemplate, a)
}
