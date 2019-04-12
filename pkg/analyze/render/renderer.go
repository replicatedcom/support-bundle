package render

import (
	"context"
	"encoding/json"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	yaml "gopkg.in/yaml.v2"
)

type Encoder interface {
	Encode(v interface{}) error
}

type Renderer struct {
	W      io.Writer
	Format string
}

func New(w io.Writer, format string) *Renderer {
	return &Renderer{
		W:      w,
		Format: format,
	}
}

func (r *Renderer) RenderResults(ctx context.Context, results []api.Result) error {
	var enc Encoder
	switch r.Format {
	case "json":
		jsonEnc := json.NewEncoder(r.W)
		jsonEnc.SetIndent("", "    ")
		enc = jsonEnc
	case "yaml":
		enc = yaml.NewEncoder(r.W)
	default:
		enc = NewHumanEncoder(r.W, true, false) // TODO: colors?
	}
	return enc.Encode(results)
}
