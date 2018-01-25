package lifecycle

import (
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tasks := []types.LifecycleTask{
		{
			Message: &types.MessageOptions{Contents: "Starting support bundle collection..."},
		},
		{
			Upload: &types.UploadOptions{},
		},
		{
			Message: &types.MessageOptions{Contents: "Upload complete! Check the analyzed bundle for more information"},
		},
	}

	lf := Lifecycle{
		GenerateTimeout:    100,
		SkipPrompts:        false,
		GenerateBundlePath: "/out",
		UploadCustomerID:   "fake",
		GraphQLClient:      nil,
		FileInfo:           nil,
		BundleTasks:        nil,
		executors:          nil,
	}

	lf.Build(tasks)
	assert.Nil(t, lf.Run())

}
