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
	producer := types.BytesProducerFromStreamProducer(producers.ReadFile("/proc/loadavg"))
	task := &plans.ByteSource{
		Producer:  producer,
		Template:  loadAverageTemplate,
		Parser:    parseLoadAvg,
		RawPath:   spec.Raw,
		JSONPath:  spec.JSON,
		HumanPath: spec.Human,
	}

	return []types.Task{task}
}

func parseLoadAvg(r io.Reader) (interface{}, error) {

	// # cat /proc/loadavg
	// 0.02 0.01 0.00 4/229 5

	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(contents), " ")
	if len(parts) != 5 {
		return nil, fmt.Errorf("expected 5 values in loadavg but found %d", len(parts))
	}

	oneMin, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, err
	}
	fiveMin, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, err
	}
	tenMin, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, err
	}

	return &LoadAverage{
		MinuteOne:  oneMin,
		MinuteFive: fiveMin,
		MinuteTen:  tenMin,
	}, nil
}
