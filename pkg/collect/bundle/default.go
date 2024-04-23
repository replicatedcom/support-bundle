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

	defaultReplicatedSpec, err := spec.Parse(defaultReplicatedYml)
	if err != nil {
		return nil, errors.Wrap(err, "parse replicated spec")
	}
	defaultSpec = append(defaultSpec, defaultReplicatedSpec...)

	return defaultSpec, nil
}

func SupportBundleSpecs() []types.Spec {
	return []types.Spec{
		{
			SpecShared: types.SpecShared{
				Description: "Support bundle version information",
				OutputDir:   "/",
			},
			SupportBundleVersion: &types.SupportBundleVersionOptions{},
		},
	}
}

// Deprecated
func CustomerJSONSpec(customerID string) types.Spec {
	return types.Spec{
		SpecShared: types.SpecShared{
			Description: "Troubleshoot Customer Metadata",
			OutputDir:   "/",
		},
		CustomerMeta: &types.CustomerMetaOptions{
			CustomerID: customerID,
		},
	}
}

func ChannelJSONSpec(channelID string) types.Spec {
	return types.Spec{
		SpecShared: types.SpecShared{
			Description: "Troubleshoot Channel Metadata",
			OutputDir:   "/",
		},
		ChannelMeta: &types.ChannelMetaOptions{
			ChannelID: channelID,
		},
	}
}

func WatchJSONSpec(watchID string) types.Spec {
	return types.Spec{
		SpecShared: types.SpecShared{
			Description: "Troubleshoot Watch Metadata",
			OutputDir:   "/",
		},
		WatchMeta: &types.WatchMetaOptions{
			WatchID: watchID,
		},
	}
}
