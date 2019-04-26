package render

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
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

func (r *Renderer) RenderBundles(ctx context.Context, bundles map[string][]collecttypes.Result, quiet bool) error {
	if quiet {
		for path := range bundles {
			fmt.Fprintf(r.W, "path: %s\n", path)
		}
		return nil
	}
	var enc Encoder
	switch r.Format {
	// "human" unsupported
	case "yaml":
		enc = yaml.NewEncoder(r.W)
	default:
		jsonEnc := json.NewEncoder(r.W)
		jsonEnc.SetIndent("", "    ")
		enc = jsonEnc
	}
	return enc.Encode(bundles)
}
