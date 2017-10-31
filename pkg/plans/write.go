package plans

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func write(rootDir string, path string, data []byte, result *types.Result) {
	dest := filepath.Join(rootDir, path)
	if err := os.MkdirAll(filepath.Dir(dest), 0744); err != nil {
		result.Error = err
		return
	}
	if err := ioutil.WriteFile(dest, data, 0644); err != nil {
		result.Error = err
	} else {
		result.Path = path
	}
}

func writeJSON(rootDir string, path string, data interface{}, result *types.Result) {
	jsonPath := filepath.Join(rootDir, path)
	if err := os.MkdirAll(filepath.Dir(jsonPath), 0744); err != nil {
		result.Error = err
		return
	}
	marshaled, err := json.Marshal(data)
	if err != nil {
		result.Error = err
		return
	}
	if err := ioutil.WriteFile(jsonPath, marshaled, 0644); err != nil {
		result.Error = err
		return
	}
	result.Path = path
}

func writeTemplate(rootDir string, path string, tmpl string, data interface{}, result *types.Result) {
	dest := filepath.Join(rootDir, path)
	if err := os.MkdirAll(filepath.Dir(dest), 0744); err != nil {
		result.Error = err
		return
	}
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		result.Error = err
		return
	}
	f, err := os.Create(dest)
	if err != nil {
		result.Error = err
		return
	}
	defer closeLogErr(f)
	if err := t.Execute(f, data); err != nil {
		result.Error = err
		return
	}
	result.Path = path
}

func writeYAML(rootDir string, path string, data interface{}, result *types.Result) {
	marshaled, err := yaml.Marshal(data)
	if err != nil {
		result.Error = err
		return
	}
	dest := filepath.Join(rootDir, path)
	if err := os.MkdirAll(filepath.Dir(dest), 0744); err != nil {
		result.Error = err
		return
	}
	if err := ioutil.WriteFile(dest, marshaled, 0644); err != nil {
		result.Error = err
		return
	}
	result.Path = path
}

// TODO context interruptible
func ioCopyContext(ctx context.Context, rootDir string, path string, data io.Reader, result *types.Result) {
	dest := filepath.Join(rootDir, path)
	if err := os.MkdirAll(filepath.Dir(dest), 0744); err != nil {
		result.Error = err
		return
	}
	f, err := os.Create(dest)
	if err != nil {
		result.Error = err
		return
	}
	n, err := io.Copy(f, data)
	if err != nil {
		result.Error = err
	}
	if n != 0 {
		result.Path = path
	}
}
