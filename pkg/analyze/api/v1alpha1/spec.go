package v1alpha1

import (
	"bytes"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

func mustAsset(name, path string) *template.Template {
	data, err := Asset(path)
	if err != nil {
		panic(err)
	}
	return template.Must(template.New(name).Delims("{{{", "}}}").Parse(string(data)))
}

func unmarshalSpec(tmpl *template.Template, data interface{}) (RawSpec, error) {
	var rawSpec RawSpec
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return rawSpec, err
	}
	err = yaml.Unmarshal(buf.Bytes(), &rawSpec)
	return rawSpec, err
}
