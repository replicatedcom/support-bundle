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

type LoadAverage struct {
	MinuteOne           float64
	MinuteFive          float64
	MinuteTen           float64
	ProcessCountRunning int
	ProcessCountTotal   int
}

const loadAverageTemplate = `
 1 Minute: {{ .MinuteOne }}
 5 Minute: {{ .MinuteFive }}
10 Minute: {{ .MinuteTen }}

{{ with .ProcessCountRunning }}Running Processes: {{ .}}{{ end }}
{{ with .ProcessCountTotal }}Total Processes: {{ . }}{{ end }}`

func PlanLoadAverage(spec types.Spec) []types.Task {
	task := &plans.ByteSource{
		Producer:  producers.ReadFile("/proc/loadavg"),
		Template:  loadAverageTemplate,
		Parser:    parseLoadAvg,
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}

func parseLoadAvg(contents []byte) (interface{}, error) {

	// # cat /proc/loadavg
	// 0.02 0.01 0.00 4/229 5

	parts := strings.Split(string(contents), " ")
	if len(parts) != 5 {
		err := fmt.Errorf("Expected 5 values in loadavg but found %d", len(parts))
		jww.ERROR.Print(err)
		return nil, err
	}

	oneMin, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		jww.ERROR.Print(err)
		return nil, err
	}
	fiveMin, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		jww.ERROR.Print(err)
		return nil, err
	}
	tenMin, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		jww.ERROR.Print(err)
		return nil, err
	}

	return &LoadAverage{
		MinuteOne:  oneMin,
		MinuteFive: fiveMin,
		MinuteTen:  tenMin,
	}, nil
}
