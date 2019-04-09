package v1

import (
	"errors"
	"text/template"
)

var (
	FileMatchesRequirementSpecTemplate *template.Template
)

func init() {
	FileMatchesRequirementSpecTemplate = mustAsset(
		"KubernetesVersionRequirementSpec", "requirements/file_matches.yml")
}

type FileMatchesRequirement struct {
	Regex   string `json:"regex" yaml:"regex" hcl:"regex"`
	Message string `json:"message" yaml:"message" hcl:"message"`
}

func (a *FileMatchesRequirement) Validate(spec AnalyzerSpec) error {
	if len(spec.CollectRefs) == 0 {
		return errors.New("collect_refs required")
	}
	if a.Regex == "" {
		return errors.New("regex required")
	}
	return nil
}

func (a *FileMatchesRequirement) GetRawSpec() (RawSpec, error) {
	return unmarshalSpec(FileMatchesRequirementSpecTemplate, a)
}
