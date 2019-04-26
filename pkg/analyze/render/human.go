package render

import (
	"fmt"
	"io"

	"github.com/mitchellh/cli"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	"github.com/replicatedcom/support-bundle/pkg/ui"
)

type HumanEncoder struct {
	UI cli.Ui
}

func NewHumanEncoder(w io.Writer, colored, force bool) *HumanEncoder {
	wUI := ui.New(nil, w, w)
	if colored {
		wUI = ui.Colored(wUI, force)
	}
	return &HumanEncoder{wUI}
}

func (e *HumanEncoder) Encode(v interface{}) error {
	switch typed := v.(type) {
	case []api.Result:
		return e.encodeResults(typed)
	default:
		return fmt.Errorf("unexpected type %T", v)
	}
}

func (e *HumanEncoder) encodeResults(results []api.Result) error {
	if len(results) == 0 {
		symbol := e.severitySymbol(common.SeverityWarn)
		e.UI.Warn(fmt.Sprintf("%s Results empty\n", symbol))
		return nil
	}

	for _, result := range results {
		if result.Message == nil {
			continue
		}
		var msg string
		symbol := e.severitySymbol(result.Severity)
		if result.Error != "" {
			msg += fmt.Sprintf(
				"%s %s\n",
				symbol, result.Error)
		} else {
			msg += fmt.Sprintf(
				"%s %s\n",
				symbol, result.Message.Primary)
		}
		msg += fmt.Sprintf(
			"%s\n",
			result.Message.Detail)
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
	case common.SeverityDebug:
		return "ℹ"
	case common.SeverityInfo:
		return "✓"
	default:
		return "✓"
	}
}
