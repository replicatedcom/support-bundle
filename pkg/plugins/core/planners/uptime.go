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

const uptimeTemplate = `
Total Time (seconds): {{ .TotalSeconds }}
Idle Time (seconds): {{ .IdleSeconds }}`

func Uptime(spec types.Spec) []types.Task {
	task := plans.StreamsSource{
		Spec:     spec,
		Producer: producers.ReadFile(types.CoreReadFileOptions{Filepath: "/proc/uptime"}),
		Template: uptimeTemplate,
		Parser:   parseUptime,
	}
	var err error
	task, err = plans.SetCommonFieldsStreamsSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
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
