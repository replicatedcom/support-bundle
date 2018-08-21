package v1

import (
	"errors"
	"text/template"
)

var (
	DockerVersionRequirementSpecTemplate *template.Template
)

func init() {
	DockerVersionRequirementSpecTemplate = mustAsset(
		"DockerVersionRequirementSpec", "requirements/docker_version.yml")
}

type DockerVersionRequirement struct {
	ServerVersionMin string `json:"server_version_minimum" yaml:"server_version_minimum" hcl:"server_version_minimum"`
}

func (a *DockerVersionRequirement) Validate(spec AnalyzerSpec) error {
	if len(spec.CollectRefs) == 0 {
		return errors.New("collect_refs required")
	}
	if a.ServerVersionMin == "" {
		return errors.New("server_version_minimum required")
	}
	return nil
}

func (a *DockerVersionRequirement) GetRawSpec() (RawSpec, error) {
	return unmarshalSpec(DockerVersionRequirementSpecTemplate, a)
}
