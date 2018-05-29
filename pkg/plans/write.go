package plans

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type writerFn func(ctx context.Context, dst io.Writer) (int64, error)

func writeResult(ctx context.Context, rootDir, path string, result *types.Result, src io.Reader) {
	writeResultFile(ctx, rootDir, path, result, func(ctx context.Context, dst io.Writer) (int64, error) {
		return write(ctx, dst, src)
	})
}

func writeResultBytes(ctx context.Context, rootDir, path string, result *types.Result, src []byte) {
	writeResultFile(ctx, rootDir, path, result, func(ctx context.Context, dst io.Writer) (int64, error) {
		return writeBytes(ctx, dst, src)
	})
}

func writeResultJSON(ctx context.Context, rootDir, path string, result *types.Result, src interface{}) {
	writeResultFile(ctx, rootDir, path, result, func(ctx context.Context, dst io.Writer) (int64, error) {
		return writeJSON(ctx, dst, src)
	})
}

func writeResultYAML(ctx context.Context, rootDir, path string, result *types.Result, src interface{}) {
	writeResultFile(ctx, rootDir, path, result, func(ctx context.Context, dst io.Writer) (int64, error) {
		return writeYAML(ctx, dst, src)
	})
}

func writeResultTemplate(ctx context.Context, rootDir, path string, result *types.Result, tmpl string, src interface{}) {
	writeResultFile(ctx, rootDir, path, result, func(ctx context.Context, dst io.Writer) (int64, error) {
		return writeTemplate(ctx, dst, tmpl, src)
	})
}

func writeResultFile(ctx context.Context, rootDir, path string, result *types.Result, fn writerFn) {
	filename := filepath.Join(rootDir, path)
	dst, err := openFile(filename)
	if err != nil {
		result.Error = err
		return
	}

	size, err := fn(ctx, dst)
	result.Error = err
	result.Path = path
	result.Size = size
}

func write(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	return ioCopyContext(ctx, dst, src)
}

func writeBytes(ctx context.Context, dst io.Writer, src []byte) (int64, error) {
	return write(ctx, dst, bytes.NewBuffer(src))
}

func writeJSON(ctx context.Context, dst io.Writer, src interface{}) (int64, error) {
	pr, pw := io.Pipe()
	defer pr.Close()
	enc := json.NewEncoder(pw)
	enc.SetIndent("", "  ")
	go func() {
		pw.CloseWithError(enc.Encode(src))
	}()
	return write(ctx, dst, pr)
}

func writeYAML(ctx context.Context, dst io.Writer, src interface{}) (int64, error) {
	marshaled, err := yaml.Marshal(src)
	if err != nil {
		return 0, err
	}
	return writeBytes(ctx, dst, marshaled)
}

func writeTemplate(ctx context.Context, dst io.Writer, tmpl string, src interface{}) (int64, error) {
	pr, pw := io.Pipe()
	defer pr.Close()
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		return 0, err
	}
	go func() {
		pw.CloseWithError(t.Execute(pw, src))
	}()
	return write(ctx, dst, pr)
}

func openFile(filename string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0744); err != nil {
		return nil, err
	}
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

// TODO context interruptible
func ioCopyContext(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}
