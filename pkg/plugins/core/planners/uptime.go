package planners

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"

	jww "github.com/spf13/jwalterweatherman"
)

type uptime struct {
	TotalSeconds float64 `json:"total_seconds"`
	IdleSeconds  float64 `json:"idle_seconds"`
}

const uptimeTmpl = `
Total Time (seconds): {{ .TotalSeconds }}
Idle Time (seconds): {{ .IdleSeconds }}`

func Uptime(spec types.Spec) []types.Task {
	task := &plans.ByteSource{
		Producer: producers.ReadFile("/proc/uptime"),
		Parser:   parseUptime,
		Template: uptimeTmpl,

		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}

func parseUptime(contents []byte) (interface{}, error) {

	// # cat /proc/uptime
	// 33524.72 66785.42

	parts := strings.Split(strings.TrimSpace(string(contents)), " ")
	if len(parts) != 2 {
		err := fmt.Errorf("Expected 2 values in uptime but found %d", len(parts))
		jww.ERROR.Print(err)
		return nil, err
	}

	totalSeconds, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		jww.ERROR.Print(err)
		return nil, err
	}
	idleSeconds, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		jww.ERROR.Print(err)
		return nil, err
	}

	return uptime{
		TotalSeconds: totalSeconds,
		IdleSeconds:  idleSeconds,
	}, nil
}
