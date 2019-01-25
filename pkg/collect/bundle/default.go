package bundle

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/bundle/defaultspec"
	"github.com/replicatedcom/support-bundle/pkg/collect/spec"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	defaultCoreYml       = defaultspec.MustAsset("assets/core.yml")
	defaultDockerYml     = defaultspec.MustAsset("assets/docker.yml")
	defaultKubernetesYml = defaultspec.MustAsset("assets/kubernetes.yml")
	defaultReplicatedYml = defaultspec.MustAsset("assets/replicated.yml")
)

func DefaultSpecs() ([]types.Spec, error) {
	var defaultSpec []types.Spec

	defaultCoreSpec, err := spec.Parse(defaultCoreYml)
	if err != nil {
		return nil, errors.Wrap(err, "parse core spec")
	}
	defaultSpec = append(defaultSpec, defaultCoreSpec...)

	defaultDockerSpec, err := spec.Parse(defaultDockerYml)
	if err != nil {
		return nil, errors.Wrap(err, "parse docker spec")
	}
	defaultSpec = append(defaultSpec, defaultDockerSpec...)

	defaultKubernetesSpec, err := spec.Parse(defaultKubernetesYml)
	if err != nil {
		return nil, errors.Wrap(err, "parse kubernetes spec")
	}
	defaultSpec = append(defaultSpec, defaultKubernetesSpec...)

	defaultReplicatedSpec, err := spec.Parse(defaultReplicatedYml)
	if err != nil {
		return nil, errors.Wrap(err, "parse replicated spec")
	}
	defaultSpec = append(defaultSpec, defaultReplicatedSpec...)

	return defaultSpec, nil
}

func SupportBundleVersionSpec() types.Spec {
	return types.Spec{
		SpecShared: types.SpecShared{
			Description: "Support bundle version information",
			OutputDir:   "/",
		},
		SupportBundleVersion: &types.SupportBundleVersionOptions{},
	}
}

func CustomerJSONSpec(customerID string) types.Spec {
	return types.Spec{
		SpecShared: types.SpecShared{
			Description: "Support Bundle Customer Metadata",
			OutputDir:   "/",
		},
		CustomerMeta: &types.CustomerMetaOptions{
			CustomerID: customerID,
		},
	}
}
