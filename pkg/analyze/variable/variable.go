package variable

import (
	"github.com/pkg/errors"
	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func Extract(v Interface, bundleReader bundlereader.BundleReader, data interface{}) (interface{}, error) {
	for _, result := range v.MatchResults(bundleReader) {
		value, err := extract(v, bundleReader, result, data)
		if err != nil || value != nil {
			return value, errors.Wrapf(err, "result %s", result.Path)
		}
	}
	return nil, nil
}

func extract(v Interface, bundleReader bundlereader.BundleReader, result collecttypes.Result, data interface{}) (interface{}, error) {
	if result.Size <= 0 {
		return nil, nil
	}
	r, err := bundleReader.Open(result.Path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	defer r.Close()
	value, err := v.ExtractValue(r, result, data)
	return value, errors.Wrap(err, "extract value")
}
