package planners

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type uptime struct {
	TotalSeconds float64 `json:"total_seconds"`
	IdleSeconds  float64 `json:"idle_seconds"`
}

const uptimeTmpl = `
Total Time (seconds): {{ .TotalSeconds }}
Idle Time (seconds): {{ .IdleSeconds }}`

func Uptime(spec types.Spec) []types.Task {
	producer := types.BytesProducerFromStreamProducer(producers.ReadFile("/proc/uptime"))
	task := &plans.ByteSource{
		Producer: producer,
		Parser:   parseUptime,
		Template: uptimeTmpl,

		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}

func parseUptime(r io.Reader) (interface{}, error) {

	// # cat /proc/uptime
	// 33524.72 66785.42

	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(strings.TrimSpace(string(contents)), " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Expected 2 values in uptime but found %d", len(parts))
	}

	totalSeconds, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, err
	}
	idleSeconds, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, err
	}

	return uptime{
		TotalSeconds: totalSeconds,
		IdleSeconds:  idleSeconds,
	}, nil
}
