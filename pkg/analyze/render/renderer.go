package render

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mitchellh/cli"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	yaml "gopkg.in/yaml.v2"
)

type Encoder interface {
	Encode(v interface{}) error
}

type Renderer struct {
	UI     cli.Ui
	Format string
}

func New(ui cli.Ui, format string) *Renderer {
	return &Renderer{
		UI:     ui,
		Format: format,
	}
}

func (r *Renderer) RenderResults(ctx context.Context, w io.Writer, results []api.Result) error {
	var enc Encoder
	switch r.Format {
	case "json":
		jsonEnc := json.NewEncoder(w)
		jsonEnc.SetIndent("", "    ")
		enc = jsonEnc
	case "yaml":
		enc = yaml.NewEncoder(w)
	default:
		enc = NewHumanEncoder(w, r.UI)
	}
	return enc.Encode(results)
}
