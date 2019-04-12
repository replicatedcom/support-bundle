package render

import (
	"fmt"
	"io"

	"github.com/mitchellh/cli"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
)

type HumanEncoder struct {
	W  io.Writer
	UI cli.Ui
}

func NewHumanEncoder(w io.Writer, ui cli.Ui) *HumanEncoder {
	return &HumanEncoder{w, ui}
}

func (e *HumanEncoder) Encode(v interface{}) error {
	results, ok := v.([]api.Result)
	if !ok {
		return fmt.Errorf("unexpected type %T", v)
	}
	return e.encodeResults(results)
}

func (e *HumanEncoder) encodeResults(results []api.Result) error {
	for _, result := range results {
		// TODO
		msg := fmt.Sprintf(
			"%s %s\n",
			e.severitySymbol(result.Severity),
			result.Message.Detail)
		if result.Message != nil {
			msg += fmt.Sprintf(
				"%s\n",
				result.Message.Primary)
		}
		e.severityOut(result.Severity)(msg)
	}
	return nil
}

func (e *HumanEncoder) severityOut(severity common.Severity) func(string) {
	switch severity {
	case common.SeverityError:
		return e.UI.Error
	case common.SeverityWarn:
		return e.UI.Warn
	case common.SeverityInfo:
		return e.UI.Info
	default:
		return e.UI.Output
	}
}

func (e *HumanEncoder) severitySymbol(severity common.Severity) string {
	switch severity {
	case common.SeverityError:
		return "✗"
	case common.SeverityWarn:
		return "!"
	case common.SeverityInfo:
		return "ℹ"
	default:
		return "✓"
	}
}
